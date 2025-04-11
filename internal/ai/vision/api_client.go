package vision

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/photoprism/photoprism/internal/api/download"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/media"
	"github.com/photoprism/photoprism/pkg/media/http/header"
	"github.com/photoprism/photoprism/pkg/media/http/scheme"
	"github.com/photoprism/photoprism/pkg/rnd"
)

// NewApiRequest returns a new Vision API request with the specified file payload and scheme.
func NewApiRequest(images Files, fileScheme string) (*ApiRequest, error) {
	imageUrls := make(Files, len(images))

	if fileScheme == scheme.Https && !strings.HasPrefix(DownloadUrl, "https://") {
		log.Tracef("vision: file request scheme changed from https to data because https is not configured")
		fileScheme = scheme.Data
	}

	for i := range images {
		switch fileScheme {
		case scheme.Https:
			if id, err := download.Register(images[i]); err != nil {
				return nil, fmt.Errorf("%s (create download url)", err)
			} else {
				imageUrls[i] = fmt.Sprintf("%s/%s", DownloadUrl, id)
			}
		case scheme.Data:
			if file, err := os.Open(images[i]); err != nil {
				return nil, fmt.Errorf("%s (create data url)", err)
			} else {
				imageUrls[i] = media.DataUrl(file)
			}
		default:
			return nil, fmt.Errorf("invalid file scheme %s", clean.Log(fileScheme))
		}
	}

	return &ApiRequest{
		Id:     rnd.UUID(),
		Model:  "",
		Images: imageUrls,
	}, nil
}

// PerformApiRequest performs a Vision API request and returns the result.
func PerformApiRequest(apiRequest *ApiRequest, uri, method, key string) (apiResponse *ApiResponse, err error) {
	if apiRequest == nil {
		return apiResponse, errors.New("api request is nil")
	}

	data, jsonErr := apiRequest.JSON()

	if jsonErr != nil {
		return apiResponse, jsonErr
	}

	// Create HTTP client and authenticated service API request.
	client := http.Client{Timeout: ServiceTimeout}
	req, reqErr := http.NewRequest(method, uri, bytes.NewReader(data))

	// Add "application/json" content type header.
	header.SetContentType(req, header.ContentTypeJson)

	// Add an authentication header if an access token is configured.
	if key != "" {
		header.SetAuthorization(req, key)
	}

	if reqErr != nil {
		return apiResponse, reqErr
	}

	// Perform API request.
	clientResp, clientErr := client.Do(req)

	if clientErr != nil {
		return apiResponse, clientErr
	}

	apiResponse = &ApiResponse{}

	// Unmarshal response and add labels, if returned.
	if apiJson, apiErr := io.ReadAll(clientResp.Body); apiErr != nil {
		return apiResponse, apiErr
	} else if apiErr = json.Unmarshal(apiJson, apiResponse); apiErr != nil {
		return apiResponse, apiErr
	} else if clientResp.StatusCode >= 300 {
		log.Debugf("vision: %s (status code %d)", apiJson, clientResp.StatusCode)
	}

	return apiResponse, nil
}
