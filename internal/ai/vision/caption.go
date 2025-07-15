package vision

import (
	"errors"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/pkg/media"
)

// CaptionPromptDefault is the default prompt used to generate captions.
var CaptionPromptDefault = `Create a caption that sounds natural and briefly describes the main content of the image in up to` +
	` three sentences for use in a photo management application. Begin with the type or number of subjects and` +
	` action. Omit text formatting and avoid meta-language such as "this picture", "the picture", "the photo",` +
	` "there are", "here is", or "a picture of". Use explicit language to describe the scene if it helps to` +
	` properly understand the picture.`

// CaptionModelDefault is the default model used to generate captions.
var CaptionModelDefault = "qwen2.5vl"

// Caption returns generated captions for the specified images.
func Caption(images Files, src media.Src) (result *CaptionResult, model *Model, err error) {
	// Return if there is no configuration or no image classification models are configured.
	if Config == nil {
		return result, model, errors.New("vision service is not configured")
	} else if model = Config.Model(ModelTypeCaption); model != nil {
		// Use remote service API if a server endpoint has been configured.
		if uri, method := model.Endpoint(); uri != "" && method != "" {
			var apiRequest *ApiRequest
			var apiResponse *ApiResponse

			if apiRequest, err = NewApiRequest(model.EndpointRequestFormat(), images, model.EndpointFileScheme()); err != nil {
				return result, model, err
			}

			if model.Name != "" {
				apiRequest.Model = model.Name
			}

			if model.Version != "" {
				apiRequest.Version = model.Version
			} else {
				apiRequest.Version = "latest"
			}

			if model.Prompt != "" {
				apiRequest.Prompt = model.Prompt
			} else {
				apiRequest.Prompt = CaptionPromptDefault
			}

			// Log JSON request data in trace mode.
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
