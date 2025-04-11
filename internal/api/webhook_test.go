package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/api/hooks"
	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/rnd"
)

func TestPostWebhook(t *testing.T) {
	app, router, conf := NewApiTest()
	conf.SetAuthMode(config.AuthModePasswd)
	defer conf.SetAuthMode(config.AuthModePublic)
	Webhook(router)
	t.Run("Success", func(t *testing.T) {
		payload := hooks.Payload{
			Type:      "api.downloads.register",
			Timestamp: time.Now().UTC(),
			Data: event.Data{
				"uuid":     rnd.UUID(),
				"filename": fs.Abs("./testdata/cat_224x224.jpg"),
			},
		}

		body := payload.JSON()
		token := "9d8b8801ffa23eb52e08ca7766283799ddfd8dd368212123"

		t.Logf("request: %s", string(body))

		response := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/webhook/instance", string(body), token)

		assert.Equal(t, http.StatusOK, response.Code)
	})
	t.Run("InvalidData", func(t *testing.T) {
		payload := hooks.Payload{
			Type:      "api.downloads.register",
			Timestamp: time.Now().UTC(),
			Data: event.Data{
				"uuid":     12345,
				"filename": fs.Abs("./testdata/green_224x224.jpg"),
			},
		}

		body := payload.JSON()
		token := "778f0f7d80579a072836c65b786145d6e0127505194cc51e"

		t.Logf("request: %s", string(body))

		response := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/webhook/instance", string(body), token)

		assert.Equal(t, http.StatusOK, response.Code)
	})
	t.Run("Unauthorized", func(t *testing.T) {
		r := PerformRequest(app, http.MethodPost, "/api/v1/webhook/instance")
		assert.Equal(t, http.StatusUnauthorized, r.Code)
	})
}
