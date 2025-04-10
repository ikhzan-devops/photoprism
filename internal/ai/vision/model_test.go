package vision

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	t.Run("Nasnet", func(t *testing.T) {
		ServiceUri = "https://app.localssl.dev/api/v1/vision"
		uri, method := NasnetModel.Endpoint()
		ServiceUri = ""
		assert.Equal(t, "https://app.localssl.dev/api/v1/vision/labels", uri)
		assert.Equal(t, http.MethodPost, method)
		uri, method = NasnetModel.Endpoint()
		assert.Equal(t, "", uri)
		assert.Equal(t, "", method)
	})
}
