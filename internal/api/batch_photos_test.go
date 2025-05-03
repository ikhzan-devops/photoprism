package api

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/photoprism/photoprism/pkg/i18n"
)

func TestBatchPhotos(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		app, router, _ := NewApiTest()

		BatchPhotos(router)

		response := PerformRequestWithBody(app,
			"POST", "/api/v1/batch/photos",
			`{"photos": ["ps6sg6be2lvl0yh7", "ps6sg6be2lvl0yh8", "ps6sg6be2lvl0ycc"]}`,
		)

		body := response.Body.String()

		assert.NotEmpty(t, body)
		assert.True(t, strings.HasPrefix(body, `{"photos":[{"ID"`), "unexpected response")

		// fmt.Println(body)
		/* photos := gjson.Get(body, "photos")
		values := gjson.Get(body, "values")
		t.Logf("photos: %#v", photos)
		t.Logf("values: %#v", values) */

		assert.Equal(t, http.StatusOK, response.Code)
	})
	t.Run("MissingSelection", func(t *testing.T) {
		app, router, _ := NewApiTest()
		BatchPhotos(router)
		r := PerformRequestWithBody(app, "POST", "/api/v1/batch/photos", `{"photos": []}`)
		val := gjson.Get(r.Body.String(), "error")
		assert.Equal(t, i18n.Msg(i18n.ErrNoItemsSelected), val.String())
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		app, router, _ := NewApiTest()
		BatchPhotos(router)
		r := PerformRequestWithBody(app, "POST", "/api/v1/batch/photos", `{"photos": 123}`)
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}
