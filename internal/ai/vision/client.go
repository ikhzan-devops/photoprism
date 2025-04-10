package vision

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/photoprism/photoprism/pkg/media/http/header"
)

// PerformApiRequest performs a Vision API request and returns the result.
func PerformApiRequest(apiRequest *ApiRequest, uri, method, key string) (apiResponse *ApiResponse, err error) {
	if apiRequest == nil {
		return apiResponse, errors.New("api request is nil")
	}

	data, jsonErr := apiRequest.MarshalJSON()

	if jsonErr != nil {
		return apiResponse, jsonErr
	}

	// Create HTTP client and authenticated service API request.
	client := http.Client{Timeout: ServiceTimeout}
	req, reqErr := http.NewRequest(method, uri, bytes.NewReader(data))

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
	}

	return apiResponse, nil
}
