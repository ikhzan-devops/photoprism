package vision

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/photoprism/photoprism/internal/ai/face"
	"github.com/photoprism/photoprism/internal/thumb/crop"
	"github.com/photoprism/photoprism/pkg/media/http/header"
	"github.com/photoprism/photoprism/pkg/media/http/scheme"
)

// Faces runs face detection and facenet algorithms over the provided source image.
func Faces(fileName string, minSize int, cacheCrop bool, expected int) (faces face.Faces, err error) {
	if fileName == "" {
		return faces, errors.New("missing image filename")
	}

	// Return if there is no configuration or no image classification models are configured.
	if Config == nil {
		return faces, errors.New("vision service is not configured")
	} else if model := Config.Model(ModelTypeFaceEmbeddings); model != nil {
		faces, err = face.Detect(fileName, false, minSize)

		if err != nil {
			return faces, err
		}

		// Skip embeddings?
		if c := len(faces); c == 0 || expected > 0 && c == expected {
			return faces, nil
		}

		if uri, method := model.Endpoint(); uri != "" && method != "" {
			faceCrops := make([]string, len(faces))

			for i, f := range faces {
				if f.Area.Col == 0 && f.Area.Row == 0 {
					faceCrops[i] = ""
					continue
				}

				if _, faceCrop, imgErr := crop.ImageFromThumb(fileName, f.CropArea(), face.CropSize, cacheCrop); imgErr != nil {
					log.Errorf("faces: failed to decode image: %s", imgErr)
					faceCrops[i] = ""
				} else if faceCrop != "" {
					faceCrops[i] = faceCrop
				}
			}

			apiRequest, apiRequestErr := NewClientRequest(faceCrops, scheme.Data)

			if apiRequestErr != nil {
				return faces, apiRequestErr
			}

			if model.Name != "" {
				apiRequest.Model = model.Name
			}

			data, jsonErr := apiRequest.MarshalJSON()

			if jsonErr != nil {
				return faces, jsonErr
			}

			// Create HTTP client and authenticated service API request.
			client := http.Client{}
			req, reqErr := http.NewRequest(method, uri, bytes.NewReader(data))
			header.SetAuthorization(req, model.EndpointKey())

			if reqErr != nil {
				return faces, reqErr
			}

			// Perform API request.
			clientResp, clientErr := client.Do(req)

			if clientErr != nil {
				return faces, clientErr
			}

			apiResponse := &ApiResponse{}

			if apiJson, apiErr := io.ReadAll(clientResp.Body); apiErr != nil {
				return faces, apiErr
			} else if apiErr = json.Unmarshal(apiJson, apiResponse); apiErr != nil {
				return faces, apiErr
			}

			for i := range faces {
				if len(apiResponse.Result.Embeddings) > i {
					faces[i].Embeddings = apiResponse.Result.Embeddings[i]
				}
			}
		} else if tf := model.FaceModel(); tf != nil {
			for i, f := range faces {
				if f.Area.Col == 0 && f.Area.Row == 0 {
					continue
				}

				if img, _, imgErr := crop.ImageFromThumb(fileName, f.CropArea(), face.CropSize, cacheCrop); imgErr != nil {
					log.Errorf("faces: failed to decode image: %s", imgErr)
				} else if embeddings := tf.Run(img); !embeddings.Empty() {
					faces[i].Embeddings = embeddings
				}
			}
		} else {
			return faces, errors.New("invalid face model configuration")
		}
	} else {
		return faces, errors.New("missing face model")
	}

	return faces, nil
}
