package vision

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/media/http/scheme"
)

func TestNewClientRequest(t *testing.T) {
	var assetsPath = fs.Abs("../../../assets")
	var examplesPath = assetsPath + "/examples"

	t.Run("Data", func(t *testing.T) {
		thumbnails := Files{examplesPath + "/chameleon_lime.jpg"}
		result, err := NewClientRequest(thumbnails, scheme.Data)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		// t.Logf("request: %#v", result)

		if result != nil {
			json, jsonErr := result.MarshalJSON()
			assert.NoError(t, jsonErr)
			assert.NotEmpty(t, json)
			// t.Logf("json: %s", json)
		}
	})

	t.Run("Https", func(t *testing.T) {
		thumbnails := Files{examplesPath + "/chameleon_lime.jpg"}
		result, err := NewClientRequest(thumbnails, scheme.Https)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		// t.Logf("request: %#v", result)
		if result != nil {
			json, jsonErr := result.MarshalJSON()
			assert.NoError(t, jsonErr)
			assert.NotEmpty(t, json)
			t.Logf("json: %s", json)
		}
	})
}
