package vision

import (
	"context"
	"errors"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/pkg/media"
)

// Caption returns generated captions for the specified images.
func Caption(images Files, mediaSrc media.Src) (result *CaptionResult, model *Model, err error) {
	// Return if there is no configuration or no image classification models are configured.
	if Config == nil {
		return result, model, errors.New("vision service is not configured")
	} else if model = Config.Model(ModelTypeCaption); model != nil {
		// Use remote service API if a server endpoint has been configured.
		if uri, method := model.Endpoint(); uri != "" && method != "" {
			var apiRequest *ApiRequest
			var apiResponse *ApiResponse

			if provider, ok := ProviderFor(model.EndpointRequestFormat()); ok && provider.Builder != nil {
				if apiRequest, err = provider.Builder.Build(context.Background(), model, images); err != nil {
					return result, model, err
				}
			} else if apiRequest, err = NewApiRequest(model.EndpointRequestFormat(), images, model.EndpointFileScheme()); err != nil {
				return result, model, err
			}

			if apiRequest.Model == "" {
				switch model.Service.RequestFormat {
				case ApiFormatOllama:
					apiRequest.Model, _, _ = model.Model()
				default:
					_, apiRequest.Model, apiRequest.Version = model.Model()
				}
			}

			apiRequest.System = model.GetSystemPrompt()
			apiRequest.Prompt = model.GetPrompt()
			apiRequest.Options = model.GetOptions()
			apiRequest.WriteLog()

			if apiResponse, err = PerformApiRequest(apiRequest, uri, method, model.EndpointKey()); err != nil {
				return result, model, err
			} else if apiResponse.Result.Caption == nil {
				return result, model, errors.New("invalid caption model response")
			}

			// Set image as the default caption source.
			if apiResponse.Result.Caption.Text != "" && apiResponse.Result.Caption.Source == "" {
				apiResponse.Result.Caption.Source = entity.SrcImage
			}

			result = apiResponse.Result.Caption
		} else {
			return result, model, errors.New("invalid caption model configuration")
		}
	} else {
		return result, model, errors.New("missing caption model")
	}

	return result, model, nil
}
