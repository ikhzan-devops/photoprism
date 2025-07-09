package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestGetPlacesReverse(t *testing.T) {
	t.Run("PublicMode", func(t *testing.T) {
		app, router, _ := NewApiTest()

		GetPlacesReverse(router)

		r := PerformRequest(app, http.MethodGet, "/api/v1/places/reverse?lat=52.5108869&lng=13.398947")

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Equal(t, "Berlin", gjson.Get(r.Body.String(), "name").String())
	})
}

func TestGetPlacesSearch(t *testing.T) {
	t.Run("PublicMode", func(t *testing.T) {
		app, router, _ := NewApiTest()

		GetPlacesSearch(router)

		r := PerformRequest(app, http.MethodGet, "/api/v1/places/search?q=Berlin&locale=en")

		assert.Equal(t, http.StatusOK, r.Code)
		assert.LessOrEqual(t, 1, int(gjson.Get(r.Body.String(), "#").Int()))
	})
}
