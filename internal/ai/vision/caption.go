package vision

import (
	"errors"
	"fmt"
	"net/url"
	"slices"

	"github.com/photoprism/photoprism/internal/api/download"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/media"
	"github.com/photoprism/photoprism/pkg/media/http/scheme"
	"github.com/photoprism/photoprism/pkg/rnd"
)

// Caption returns generated captions for the specified images.
func Caption(imgName string, src media.Src) (result CaptionResult, err error) {
	// Return if there is no configuration or no image classification models are configured.
	if Config == nil {
		return result, errors.New("vision service is not configured")
	} else if model := Config.Model(ModelTypeCaption); model != nil {
		// Use remote service API if a server endpoint has been configured.
		if uri, method := model.Endpoint(); uri != "" && method != "" {
			var imgUrl string

			switch src {
			case media.SrcLocal:
				// Return if no thumbnail filenames were given.
				if !fs.FileExistsNotEmpty(imgName) {
					return result, errors.New("invalid image file name")
				}

				/* TODO: Add support for data URLs to the service.
				if file, fileErr := os.Open(imgName); fileErr != nil {
					return result, fmt.Errorf("%s (open image file)", err)
				} else {
					imgUrl = media.DataUrl(file)
				} */

				dlId, dlErr := download.Register(imgName)

				if dlErr != nil {
					return result, fmt.Errorf("%s (create download url)", err)
				}

				imgUrl = fmt.Sprintf("%s/%s", DownloadUrl, dlId)
			case media.SrcRemote:
				var u *url.URL
				if u, err = url.Parse(imgName); err != nil {
					return result, fmt.Errorf("%s (invalid image url)", err)
				} else if !slices.Contains(scheme.HttpsHttp, u.Scheme) {
					return result, fmt.Errorf("unsupported image url scheme %s", clean.Log(u.Scheme))
				} else {
					imgUrl = u.String()
				}
			default:
				return result, fmt.Errorf("unsupported media source type %s", clean.Log(src))
			}

			apiRequest := &ApiRequest{
				Id:    rnd.UUID(),
				Model: model.Name,
				Url:   imgUrl,
			}

			/* if json, _ := apiRequest.JSON(); len(json) > 0 {
				log.Debugf("request: %s", json)
			} */

			apiResponse, apiErr := PerformApiRequest(apiRequest, uri, method, model.EndpointKey())

			if apiErr != nil {
				return result, apiErr
			} else if apiResponse.Result.Caption == nil {
				return result, errors.New("invalid caption model response")
			}

			// Set image as the default caption source.
			if apiResponse.Result.Caption.Text != "" && apiResponse.Result.Caption.Source == "" {
				apiResponse.Result.Caption.Source = entity.SrcImage
			}

			result = *apiResponse.Result.Caption
		} else {
			return result, errors.New("invalid caption model configuration")
		}
	} else {
		return result, errors.New("missing caption model")
	}

	return result, nil
}
