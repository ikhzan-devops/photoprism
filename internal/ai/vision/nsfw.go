package vision

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/photoprism/photoprism/internal/ai/nsfw"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/media"
	"github.com/photoprism/photoprism/pkg/media/http/header"
	"github.com/photoprism/photoprism/pkg/media/http/scheme"
)

// Nsfw checks the specified images for inappropriate content.
func Nsfw(images Files, src media.Src) (result []nsfw.Result, err error) {
	// Return if no thumbnail filenames were given.
	if len(images) == 0 {
		return result, errors.New("at least one image required")
	}

	result = make([]nsfw.Result, len(images))

	// Return if there is no configuration or no image classification models are configured.
	if Config == nil {
		return result, errors.New("vision service is not configured")
	} else if model := Config.Model(ModelTypeNsfw); model != nil {
		// Use remote service API if a server endpoint has been configured.
		if uri, method := model.Endpoint(); uri != "" && method != "" {
			apiRequest, apiRequestErr := NewClientRequest(images, scheme.Data)

			if apiRequestErr != nil {
				return result, apiRequestErr
			}

			if model.Name != "" {
				apiRequest.Model = model.Name
			}

			data, jsonErr := apiRequest.MarshalJSON()

			if jsonErr != nil {
				return result, jsonErr
			}

			// Create HTTP client and authenticated service API request.
			client := http.Client{}
			req, reqErr := http.NewRequest(method, uri, bytes.NewReader(data))
			header.SetAuthorization(req, model.EndpointKey())

			if reqErr != nil {
				return result, reqErr
			}

			// Perform API request.
			clientResp, clientErr := client.Do(req)

			if clientErr != nil {
				return result, clientErr
			}

			apiResponse := &ApiResponse{}

			// Unmarshal response and add labels, if returned.
			if apiJson, apiErr := io.ReadAll(clientResp.Body); apiErr != nil {
				return result, apiErr
			} else if apiErr = json.Unmarshal(apiJson, apiResponse); apiErr != nil {
				return result, apiErr
			}

			result = apiResponse.Result.Nsfw
		} else if tf := model.NsfwModel(); tf != nil {
			// Predict labels with local TensorFlow model.
			for i := range images {
				var labels nsfw.Result

				switch src {
				case media.SrcLocal:
					labels, err = tf.File(images[i])
				case media.SrcRemote:
					labels, err = tf.Url(images[i])
				default:
					return result, fmt.Errorf("invalid image source %s", clean.Log(src))
				}

				if err != nil {
					log.Errorf("nsfw: %s", err)
				}

				result[i] = labels
			}
		} else {
			return result, errors.New("invalid nsfw model configuration")
		}
	} else {
		return result, errors.New("missing nsfw model")
	}

	return result, nil
}
