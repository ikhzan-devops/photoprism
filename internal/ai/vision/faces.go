package vision

import (
	"errors"

	"github.com/photoprism/photoprism/internal/ai/face"
	"github.com/photoprism/photoprism/internal/thumb/crop"
	"github.com/photoprism/photoprism/pkg/media/http/scheme"
)

// Faces runs face detection and facenet algorithms over the provided source image.
func Faces(fileName string, minSize int, cacheCrop bool, expected int) (result face.Faces, err error) {
	if fileName == "" {
		return result, errors.New("missing image filename")
	}

	// Return if there is no configuration or no image classification models are configured.
	if Config == nil {
		return result, errors.New("vision service is not configured")
	} else if model := Config.Model(ModelTypeFaceEmbeddings); model != nil {
		result, err = face.Detect(fileName, false, minSize)

		if err != nil {
			return result, err
		}

		// Skip embeddings?
		if c := len(result); c == 0 || expected > 0 && c == expected {
			return result, nil
		}

		if uri, method := model.Endpoint(); uri != "" && method != "" {
			faceCrops := make([]string, len(result))

			for i, f := range result {
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
				return result, apiRequestErr
			}

			if model.Name != "" {
				apiRequest.Model = model.Name
			}

			apiResponse, apiErr := PerformApiRequest(apiRequest, uri, method, model.EndpointKey())

			if apiErr != nil {
				return result, apiErr
			}

			for i := range result {
				if len(apiResponse.Result.Embeddings) > i {
					result[i].Embeddings = apiResponse.Result.Embeddings[i]
				}
			}
		} else if tf := model.FaceModel(); tf != nil {
			for i, f := range result {
				if f.Area.Col == 0 && f.Area.Row == 0 {
					continue
				}

				if img, _, imgErr := crop.ImageFromThumb(fileName, f.CropArea(), face.CropSize, cacheCrop); imgErr != nil {
					log.Errorf("faces: failed to decode image: %s", imgErr)
				} else if embeddings := tf.Run(img); !embeddings.Empty() {
					result[i].Embeddings = embeddings
				}
			}
		} else {
			return result, errors.New("invalid face model configuration")
		}
	} else {
		return result, errors.New("missing face model")
	}

	return result, nil
}
