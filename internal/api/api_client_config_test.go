package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/photoprism/photoprism/internal/config"
)

func TestGetClientConfig(t *testing.T) {
	t.Run("Public", func(t *testing.T) {
		app, router, _ := NewApiTest()
		GetClientConfig(router)
		r := PerformRequest(app, "GET", "/api/v1/config")
		val := gjson.Get(r.Body.String(), "mode")
		assert.Equal(t, http.StatusOK, r.Code)
		assert.Equal(t, "user", val.String())
	})
	t.Run("Unauthorized", func(t *testing.T) {
		app, router, conf := NewApiTest()
		conf.SetAuthMode(config.AuthModePasswd)
		defer conf.SetAuthMode(config.AuthModePublic)
		GetClientConfig(router)
		r := AuthenticatedRequest(app, "GET", "/api/v1/config", "")
		val := gjson.Get(r.Body.String(), "mode")
		assert.Equal(t, http.StatusOK, r.Code)
		assert.Equal(t, "public", val.String())
	})
	t.Run("FrontendDisabled", func(t *testing.T) {
		app, router, conf := NewApiTest()
		conf.SetAuthMode(config.AuthModePasswd)
		defer conf.SetAuthMode(config.AuthModePublic)
		conf.Options().DisableFrontend = true
		GetClientConfig(router)
		r := PerformRequest(app, "GET", "/api/v1/config")
		assert.Equal(t, http.StatusUnauthorized, r.Code)
		conf.Options().DisableFrontend = false
	})
}
