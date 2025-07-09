package api

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/pkg/i18n"
)

func TestBatchPhotosEdit(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Create new API test instance.
		app, router, _ := NewApiTest()

		// Attach POST /api/v1/batch/photos/edit request handler.
		BatchPhotosEdit(router)

		// Specify the unique IDs of the photos used for testing.
		photoUIDs := `["ps6sg6be2lvl0yh7", "ps6sg6be2lvl0yh8"]`

		// Get the photo models and current values for the batch edit form.
		editResponse := PerformRequestWithBody(app,
			"POST", "/api/v1/batch/photos/edit",
			fmt.Sprintf(`{"photos": %s}`, photoUIDs),
		)

		// Check the edit response status code.
		assert.Equal(t, http.StatusOK, editResponse.Code)

		// Check the edit response body.
		editBody := editResponse.Body.String()
		assert.NotEmpty(t, editBody)
		assert.True(t, strings.HasPrefix(editBody, `{"models":[{"ID"`), "unexpected response")

		// Check the edit response values.
		editValues := gjson.Get(editBody, "values").Raw
		t.Logf("edit values: %#v", editValues)

		// Send the edit form values back to the same API endpoint and check for errors.
		saveResponse := PerformRequestWithBody(app,
			"POST", "/api/v1/batch/photos/edit",
			fmt.Sprintf(`{"photos": %s, "values": %s}`, photoUIDs, editValues),
		)

		// Check the save response status code.
		assert.Equal(t, http.StatusOK, saveResponse.Code)

		// Check the save response body.
		saveBody := saveResponse.Body.String()
		assert.NotEmpty(t, saveBody)

		// Check the save response values.
		saveValues := gjson.Get(saveBody, "values").Raw
		t.Logf("save values: %#v", saveValues)
		assert.Equal(t, editValues, saveValues)
	})
	t.Run("ReturnPhotosAndValues", func(t *testing.T) {
		app, router, conf := NewApiTest()
		conf.SetAuthMode(config.AuthModePasswd)
		defer conf.SetAuthMode(config.AuthModePublic)
		authToken := AuthenticateUser(app, router, "alice", "Alice123!")

		// Attach POST /api/v1/batch/photos/edit request handler.
		BatchPhotosEdit(router)

		response := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/batch/photos/edit",
			`{"photos": ["ps6sg6be2lvl0yh7","ps6sg6be2lvl0yh8","ps6sg6byk7wrbk47","ps6sg6be2lvl0yh0"], "return": true, "values": {}}`,
			authToken)

		body := response.Body.String()

		assert.NotEmpty(t, body)
		assert.True(t, strings.HasPrefix(body, `{"models":[{"ID"`), "unexpected response")

		fmt.Println(body)
		/* models := gjson.Get(body, "models")
		values := gjson.Get(body, "values")
		t.Logf("models: %#v", models)
		t.Logf("values: %#v", values) */

		assert.Equal(t, http.StatusOK, response.Code)
	})
	t.Run("MissingSelection", func(t *testing.T) {
		app, router, _ := NewApiTest()

		// Attach POST /api/v1/batch/photos/edit request handler.
		BatchPhotosEdit(router)

		r := PerformRequestWithBody(app, "POST", "/api/v1/batch/photos/edit", `{"photos": [], "return": true}`)
		val := gjson.Get(r.Body.String(), "error")
		assert.Equal(t, i18n.Msg(i18n.ErrNoItemsSelected), val.String())
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		app, router, _ := NewApiTest()

		// Attach POST /api/v1/batch/photos/edit request handler.
		BatchPhotosEdit(router)

		r := PerformRequestWithBody(app, "POST", "/api/v1/batch/photos/edit", `{"photos": 123, "return": true}`)
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
	t.Run("ReturnValuesAsAdmin", func(t *testing.T) {
		app, router, conf := NewApiTest()

		conf.SetAuthMode(config.AuthModePasswd)
		defer conf.SetAuthMode(config.AuthModePublic)

		// Attach POST /api/v1/batch/photos/edit request handler.
		BatchPhotosEdit(router)

		sessId := AuthenticateUser(app, router, "alice", "Alice123!")

		response := AuthenticatedRequestWithBody(app,
			"POST", "/api/v1/batch/photos/edit",
			`{"photos": ["ps6sg6be2lvl0yh7", "ps6sg6be2lvl0yh8"]}`,
			sessId,
		)

		body := response.Body.String()

		assert.NotEmpty(t, body)
		assert.True(t, strings.HasPrefix(body, `{"models":[{"ID"`), "unexpected response")

		assert.Equal(t, http.StatusOK, response.Code)
	})
	t.Run("ReturnValuesAsGuest", func(t *testing.T) {
		app, router, conf := NewApiTest()

		conf.SetAuthMode(config.AuthModePasswd)
		defer conf.SetAuthMode(config.AuthModePublic)

		// Attach POST /api/v1/batch/photos/edit request handler.
		BatchPhotosEdit(router)

		sessId := AuthenticateUser(app, router, "gandalf", "Gandalf123!")

		response := AuthenticatedRequestWithBody(app,
			"POST", "/api/v1/batch/photos/edit",
			`{"photos": ["ps6sg6be2lvl0yh7", "ps6sg6be2lvl0yh8"]}`,
			sessId,
		)

		if response.Code != http.StatusForbidden {
			t.Fatal(response.Body.String())
		}

		val := gjson.Get(response.Body.String(), "error")
		assert.Equal(t, "Permission denied", val.String())
	})
}
