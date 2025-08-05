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
	/*t.Run("SuccessNoChange", func(t *testing.T) {
		// Create new API test instance.
		app, router, _ := NewApiTest()

		// Attach POST /api/v1/batch/photos/edit request handler.
		BatchPhotosEdit(router)

		// Specify the unique IDs of the photos used for testing.
		photoUIDs := `["pqkm36fjqvset9uy", "pqkm36fjqvset9uz"]`

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
	})*/
	t.Run("SuccessChangeValues", func(t *testing.T) {
		// Create new API test instance.
		app, router, _ := NewApiTest()

		// Attach POST /api/v1/batch/photos/edit request handler.
		BatchPhotosEdit(router)

		// Specify the unique IDs of the photos used for testing.
		photoUIDs := `["pqkm36fjqvset9uy", "pqkm36fjqvset9uz"]`

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

		// Check the edit response values.
		editPhotos := gjson.Get(editBody, "models").Array()
		assert.Equal(t, len(editPhotos), 2)
		editValues := gjson.Get(editBody, "values").Raw
		timezoneBefore := gjson.Get(editValues, "TimeZone")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":true,\"action\":\"none\"}", timezoneBefore.String())
		altitudeBefore := gjson.Get(editValues, "Altitude")
		assert.Equal(t, "{\"value\":0,\"mixed\":true,\"action\":\"none\"}", altitudeBefore.String())
		typeBefore := gjson.Get(editValues, "Type")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":true,\"action\":\"none\"}", typeBefore.String())
		yearBefore := gjson.Get(editValues, "Year")
		assert.Equal(t, "{\"value\":-2,\"mixed\":true,\"action\":\"none\"}", yearBefore.String())
		dayBefore := gjson.Get(editValues, "Day")
		assert.Equal(t, "{\"value\":-2,\"mixed\":true,\"action\":\"none\"}", dayBefore.String())
		monthBefore := gjson.Get(editValues, "Month")
		assert.Equal(t, "{\"value\":-2,\"mixed\":true,\"action\":\"none\"}", monthBefore.String())
		titleBefore := gjson.Get(editValues, "Title")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":true,\"action\":\"none\"}", titleBefore.String())
		captionBefore := gjson.Get(editValues, "Caption")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":true,\"action\":\"none\"}", captionBefore.String())
		subjectBefore := gjson.Get(editValues, "DetailsSubject")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":true,\"action\":\"none\"}", subjectBefore.String())
		artistBefore := gjson.Get(editValues, "DetailsArtist")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":true,\"action\":\"none\"}", artistBefore.String())
		copyrightBefore := gjson.Get(editValues, "DetailsCopyright")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":true,\"action\":\"none\"}", copyrightBefore.String())
		licenseBefore := gjson.Get(editValues, "DetailsLicense")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":true,\"action\":\"none\"}", licenseBefore.String())
		FavoriteBefore := gjson.Get(editValues, "Favorite")
		assert.Equal(t, "{\"value\":false,\"mixed\":true,\"action\":\"none\"}", FavoriteBefore.String())
		ScanBefore := gjson.Get(editValues, "Scan")
		assert.Equal(t, "{\"value\":false,\"mixed\":true,\"action\":\"none\"}", ScanBefore.String())
		PrivateBefore := gjson.Get(editValues, "Private")
		assert.Equal(t, "{\"value\":false,\"mixed\":true,\"action\":\"none\"}", PrivateBefore.String())
		PanoramaBefore := gjson.Get(editValues, "Panorama")
		assert.Equal(t, "{\"value\":false,\"mixed\":true,\"action\":\"none\"}", PanoramaBefore.String())
		// Send the edit form values back to the same API endpoint and check for errors.
		saveResponse := PerformRequestWithBody(app,
			"POST", "/api/v1/batch/photos/edit",
			fmt.Sprintf(`{"photos": %s, "values": %s}`, photoUIDs,
				"{"+
					"\"TimeZone\":{\"value\":\"Europe/Vienna\",\"mixed\":false,\"action\":\"update\"},"+
					"\"Altitude\":{\"value\":145,\"mixed\":false,\"action\":\"update\"},"+
					"\"Year\":{\"value\":2000,\"mixed\":false,\"action\":\"update\"},"+
					"\"Month\":{\"value\":11,\"mixed\":true,\"action\":\"update\"},"+
					"\"Day\":{\"value\":-1,\"mixed\":true,\"action\":\"update\"},"+
					"\"Title\":{\"value\":\"My Batch Edited Title\",\"mixed\":false,\"action\":\"update\"},"+
					"\"Caption\":{\"value\":\"Batch edited caption\",\"mixed\":false,\"action\":\"update\"},"+
					"\"DetailsSubject\":{\"value\":\"Batch edited subject\",\"mixed\":false,\"action\":\"update\"},"+
					"\"DetailsArtist\":{\"value\":\"Batchie\",\"mixed\":false,\"action\":\"update\"},"+
					"\"DetailsCopyright\":{\"value\":\"Batch edited copyright\",\"mixed\":false,\"action\":\"update\"},"+
					"\"DetailsLicense\":{\"value\":\"Batch edited license\",\"mixed\":false,\"action\":\"update\"},"+
					"\"Type\":{\"value\":\"live\",\"mixed\":false,\"action\":\"update\"},"+
					"\"Favorite\":{\"value\":false,\"mixed\":false,\"action\":\"update\"},"+
					"\"Panorama\":{\"value\":true,\"mixed\":false,\"action\":\"update\"},"+
					"\"Private\":{\"value\":true,\"mixed\":false,\"action\":\"update\"},"+
					"\"Scan\":{\"value\":true,\"mixed\":false,\"action\":\"update\"}"+
					"}"),
		)

		// Check the save response status code.
		assert.Equal(t, http.StatusOK, saveResponse.Code)

		// Check the save response body.
		saveBody := saveResponse.Body.String()
		assert.NotEmpty(t, saveBody)

		// Check the save response values.
		saveValues := gjson.Get(saveBody, "values").Raw
		timezoneAfter := gjson.Get(saveValues, "TimeZone")
		assert.Equal(t, "{\"value\":\"Europe/Vienna\",\"mixed\":false,\"action\":\"none\"}", timezoneAfter.String())
		altitudeAfter := gjson.Get(saveValues, "Altitude")
		assert.Equal(t, "{\"value\":145,\"mixed\":false,\"action\":\"none\"}", altitudeAfter.String())
		typeAfter := gjson.Get(saveValues, "Type")
		assert.Equal(t, "{\"value\":\"live\",\"mixed\":false,\"action\":\"none\"}", typeAfter.String())
		yearAfter := gjson.Get(saveValues, "Year")
		assert.Equal(t, "{\"value\":2000,\"mixed\":false,\"action\":\"none\"}", yearAfter.String())
		dayAfter := gjson.Get(saveValues, "Day")
		assert.Equal(t, "{\"value\":-1,\"mixed\":false,\"action\":\"none\"}", dayAfter.String())
		monthAfter := gjson.Get(saveValues, "Month")
		assert.Equal(t, "{\"value\":11,\"mixed\":false,\"action\":\"none\"}", monthAfter.String())
		titleAfter := gjson.Get(saveValues, "Title")
		assert.Equal(t, "{\"value\":\"My Batch Edited Title\",\"mixed\":false,\"action\":\"none\"}", titleAfter.String())
		captionAfter := gjson.Get(saveValues, "Caption")
		assert.Equal(t, "{\"value\":\"Batch edited caption\",\"mixed\":false,\"action\":\"none\"}", captionAfter.String())
		subjectAfter := gjson.Get(saveValues, "DetailsSubject")
		assert.Equal(t, "{\"value\":\"Batch edited subject\",\"mixed\":false,\"action\":\"none\"}", subjectAfter.String())
		artistAfter := gjson.Get(saveValues, "DetailsArtist")
		assert.Equal(t, "{\"value\":\"Batchie\",\"mixed\":false,\"action\":\"none\"}", artistAfter.String())
		copyrightAfter := gjson.Get(saveValues, "DetailsCopyright")
		assert.Equal(t, "{\"value\":\"Batch edited copyright\",\"mixed\":false,\"action\":\"none\"}", copyrightAfter.String())
		licenseAfter := gjson.Get(saveValues, "DetailsLicense")
		assert.Equal(t, "{\"value\":\"Batch edited license\",\"mixed\":false,\"action\":\"none\"}", licenseAfter.String())
		FavoriteAfter := gjson.Get(saveValues, "Favorite")
		assert.Equal(t, "{\"value\":false,\"mixed\":false,\"action\":\"none\"}", FavoriteAfter.String())
		ScanAfter := gjson.Get(saveValues, "Scan")
		assert.Equal(t, "{\"value\":true,\"mixed\":false,\"action\":\"none\"}", ScanAfter.String())
		PrivateAfter := gjson.Get(saveValues, "Private")
		assert.Equal(t, "{\"value\":true,\"mixed\":false,\"action\":\"none\"}", PrivateAfter.String())
		PanoramaAfter := gjson.Get(saveValues, "Panorama")
		assert.Equal(t, "{\"value\":true,\"mixed\":false,\"action\":\"none\"}", PanoramaAfter.String())
	})
	// TODO Requires updating albums/labels functionality
	/*t.Run("SuccessChangeAlbumAndLabels", func(t *testing.T) {
		// Create new API test instance.
		app, router, _ := NewApiTest()

		// Attach POST /api/v1/batch/photos/edit request handler.
		BatchPhotosEdit(router)

		// Specify the unique IDs of the photos used for testing.
		photoUIDs := `["pqkm36fjqvset9uy", "pqkm36fjqvset9uz"]`

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

		// Check the edit response values.
		editPhotos := gjson.Get(editBody, "models").Array()
		assert.Equal(t, len(editPhotos), 2)
		editValues := gjson.Get(editBody, "values").Raw
		t.Logf(editValues)
		albumsBefore := gjson.Get(editValues, "Albums")
		assert.Equal(t, "{\"items\":[{\"value\":\"as6sg6bipotaab19\",\"title\":\"&IlikeFood\",\"mixed\":false,\"action\":\"none\"},{\"value\":\"as6sg6bxpogaaba8\",\"title\":\"Holiday 2030\",\"mixed\":true,\"action\":\"none\"},{\"value\":\"as6sg6bxpogaaba7\",\"title\":\"Christmas 2030\",\"mixed\":true,\"action\":\"none\"}],\"mixed\":true,\"action\":\"none\"}", albumsBefore.String())
		labelsBefore := gjson.Get(editValues, "Labels")
		assert.Equal(t, "{\"items\":[{\"value\":\"ls6sg6b1wowuy3c4\",\"title\":\"Cake\",\"mixed\":false,\"action\":\"none\"},{\"value\":\"ls6sg6b1wowuy3c3\",\"title\":\"Flower\",\"mixed\":true,\"action\":\"none\"},{\"value\":\"ls6sg6b1wowuy316\",\"title\":\"&friendship\",\"mixed\":true,\"action\":\"none\"}],\"mixed\":true,\"action\":\"none\"}", labelsBefore.String())

		// Send the edit form values back to the same API endpoint and check for errors.
		saveResponse := PerformRequestWithBody(app,
			"POST", "/api/v1/batch/photos/edit",
			fmt.Sprintf(`{"photos": %s, "values": %s}`, photoUIDs,
				"{"+
					"\"Labels\":{\"items\":[{\"value\":\"ls6sg6b1wowuy3c4\",\"title\":\"Cake\",\"mixed\":false,\"action\":\"remove\"},{\"value\":\"ls6sg6b1wowuy3c3\",\"title\":\"Flower\",\"mixed\":false,\"action\":\"add\"},{\"value\":\"ls6sg6b1wowuy316\",\"title\":\"&friendship\",\"mixed\":false,\"action\":\"remove\"},{\"value\":\"\",\"title\":\"BatchLabel\",\"mixed\":false,\"action\":\"add\"}],\"mixed\":false,\"action\":\"update\"},"+
					"\"Albums\":{\"items\":[{\"value\":\"as6sg6bipotaab19\",\"title\":\"&IlikeFood\",\"mixed\":false,\"action\":\"remove\"},{\"value\":\"as6sg6bxpogaaba8\",\"title\":\"Holiday 2030\",\"mixed\":true,\"action\":\"none\"},{\"value\":\"as6sg6bxpogaaba7\",\"title\":\"Christmas 2030\",\"mixed\":false,\"action\":\"add\"}, {\"value\":\"\",\"title\":\"BatchAlbum\",\"mixed\":false,\"action\":\"add\"}],\"mixed\":true,\"action\":\"update\"}"+
					"}"),
		)

		// Check the save response status code.
		assert.Equal(t, http.StatusOK, saveResponse.Code)

		// Check the save response body.
		saveBody := saveResponse.Body.String()
		assert.NotEmpty(t, saveBody)

		// Check the save response values.
		saveValues := gjson.Get(saveBody, "values").Raw
		albumsAfter := gjson.Get(saveValues, "Albums")
		assert.Equal(t, "{\"items\":[{\"value\":\"as6sg6bxpogaaba8\",\"title\":\"Holiday 2030\",\"mixed\":true,\"action\":\"none\"},{\"value\":\"as6sg6bxpogaaba7\",\"title\":\"Christmas 2030\",\"mixed\":false,\"action\":\"none\"},{\"value\":\"\",\"title\":\"BatchAlbum\",\"mixed\":false,\"action\":\"none\"}],\"mixed\":true,\"action\":\"none\"}", albumsAfter.String())
		labelsAfter := gjson.Get(saveValues, "Labels")
		assert.Equal(t, "{\"items\":[{\"value\":\"ls6sg6b1wowuy3c3\",\"title\":\"Flower\",\"mixed\":false,\"action\":\"none\"},{\"value\":\"\",\"title\":\"BatchLabel\",\"mixed\":true,\"action\":\"none\"}],\"mixed\":false,\"action\":\"none\"}", labelsAfter.String())
	})*/
	t.Run("SuccessChangeCountry", func(t *testing.T) {
		// Create new API test instance.
		app, router, _ := NewApiTest()

		// Attach POST /api/v1/batch/photos/edit request handler.
		BatchPhotosEdit(router)

		// Specify the unique IDs of the photos used for testing.
		photoUIDs := `["pqkm36fjqvset9uy", "pqkm36fjqvset9uz"]`

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

		// Check the edit response values.
		editPhotos := gjson.Get(editBody, "models").Array()
		assert.Equal(t, len(editPhotos), 2)
		editValues := gjson.Get(editBody, "values").Raw
		timezoneBefore := gjson.Get(editValues, "TimeZone")
		assert.Equal(t, "{\"value\":\"Europe/Vienna\",\"mixed\":false,\"action\":\"none\"}", timezoneBefore.String())
		altitudeBefore := gjson.Get(editValues, "Altitude")
		assert.Equal(t, "{\"value\":145,\"mixed\":false,\"action\":\"none\"}", altitudeBefore.String())
		countryBefore := gjson.Get(editValues, "Country")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":true,\"action\":\"none\"}", countryBefore.String())
		latBefore := gjson.Get(editValues, "Lat")
		assert.Equal(t, "{\"value\":0,\"mixed\":true,\"action\":\"none\"}", latBefore.String())
		lngBefore := gjson.Get(editValues, "Lng")
		assert.Equal(t, "{\"value\":0,\"mixed\":true,\"action\":\"none\"}", lngBefore.String())
		// Send the edit form values back to the same API endpoint and check for errors.
		saveResponse := PerformRequestWithBody(app,
			"POST", "/api/v1/batch/photos/edit",
			fmt.Sprintf(`{"photos": %s, "values": %s}`, photoUIDs,
				"{"+
					"\"Country\":{\"value\":\"gb\",\"mixed\":false,\"action\":\"update\"},"+
					"\"Lat\":{\"value\":0,\"mixed\":false,\"action\":\"update\"},"+
					"\"Lng\":{\"value\":0,\"mixed\":false,\"action\":\"update\"}"+
					"}"),
		)

		// Check the save response status code.
		assert.Equal(t, http.StatusOK, saveResponse.Code)

		// Check the save response body.
		saveBody := saveResponse.Body.String()
		assert.NotEmpty(t, saveBody)

		// Check the save response values.
		saveValues := gjson.Get(saveBody, "values").Raw
		timezoneAfter := gjson.Get(saveValues, "TimeZone")
		assert.Equal(t, "{\"value\":\"Europe/Vienna\",\"mixed\":false,\"action\":\"none\"}", timezoneAfter.String())
		altitudeAfter := gjson.Get(saveValues, "Altitude")
		assert.Equal(t, "{\"value\":145,\"mixed\":false,\"action\":\"none\"}", altitudeAfter.String())
		countryAfter := gjson.Get(saveValues, "Country")
		assert.Equal(t, "{\"value\":\"gb\",\"mixed\":false,\"action\":\"none\"}", countryAfter.String())
		latAfter := gjson.Get(saveValues, "Lat")
		assert.Equal(t, "{\"value\":0,\"mixed\":false,\"action\":\"none\"}", latAfter.String())
		lngAfter := gjson.Get(saveValues, "Lng")
		assert.Equal(t, "{\"value\":0,\"mixed\":false,\"action\":\"none\"}", lngAfter.String())
	})
	// TODO Requires updating timezone functionality
	/*t.Run("SuccessChangeLocationValues", func(t *testing.T) {
		// Create new API test instance.
		app, router, _ := NewApiTest()

		// Attach POST /api/v1/batch/photos/edit request handler.
		BatchPhotosEdit(router)

		// Specify the unique IDs of the photos used for testing.
		photoUIDs := `["pqkm36fjqvset9uy", "pqkm36fjqvset9uz"]`

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

		// Check the edit response values.
		editPhotos := gjson.Get(editBody, "models").Array()
		assert.Equal(t, len(editPhotos), 2)
		editValues := gjson.Get(editBody, "values").Raw
		timezoneBefore := gjson.Get(editValues, "TimeZone")
		assert.Equal(t, "{\"value\":\"Europe/Vienna\",\"mixed\":false,\"action\":\"none\"}", timezoneBefore.String())
		altitudeBefore := gjson.Get(editValues, "Altitude")
		assert.Equal(t, "{\"value\":145,\"mixed\":false,\"action\":\"none\"}", altitudeBefore.String())
		countryBefore := gjson.Get(editValues, "Country")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":true,\"action\":\"none\"}", countryBefore.String())
		latBefore := gjson.Get(editValues, "Lat")
		assert.Equal(t, "{\"value\":0,\"mixed\":true,\"action\":\"none\"}", latBefore.String())
		lngBefore := gjson.Get(editValues, "Lng")
		assert.Equal(t, "{\"value\":0,\"mixed\":true,\"action\":\"none\"}", lngBefore.String())
		// Send the edit form values back to the same API endpoint and check for errors.
		saveResponse := PerformRequestWithBody(app,
			"POST", "/api/v1/batch/photos/edit",
			fmt.Sprintf(`{"photos": %s, "values": %s}`, photoUIDs,
				"{"+
					"\"Lat\":{\"value\":21.850195,\"mixed\":false,\"action\":\"update\"},"+
					"\"Lng\":{\"value\":90.18015,\"mixed\":false,\"action\":\"update\"}"+
					"}"),
		)

		// Check the save response status code.
		assert.Equal(t, http.StatusOK, saveResponse.Code)

		// Check the save response body.
		saveBody := saveResponse.Body.String()
		assert.NotEmpty(t, saveBody)

		// Check the save response values.
		saveValues := gjson.Get(saveBody, "values").Raw
		timezoneAfter := gjson.Get(saveValues, "TimeZone")
		assert.Equal(t, "{\"value\":\"Asia/Dhaka\",\"mixed\":false,\"action\":\"none\"}", timezoneAfter.String())
		altitudeAfter := gjson.Get(saveValues, "Altitude")
		assert.Equal(t, "{\"value\":145,\"mixed\":false,\"action\":\"none\"}", altitudeAfter.String())
		countryAfter := gjson.Get(saveValues, "Country")
		assert.Equal(t, "{\"value\":\"bd\",\"mixed\":false,\"action\":\"none\"}", countryAfter.String())
		latAfter := gjson.Get(saveValues, "Lat")
		assert.Equal(t, "{\"value\":21.850195,\"mixed\":false,\"action\":\"none\"}", latAfter.String())
		lngAfter := gjson.Get(saveValues, "Lng")
		assert.Equal(t, "{\"value\":90.18015,\"mixed\":false,\"action\":\"none\"}", lngAfter.String())
	})*/
	//TODO Requires remove functionality
	/*t.Run("SuccessRemoveValues", func(t *testing.T) {
		// Create new API test instance.
		app, router, _ := NewApiTest()

		// Attach POST /api/v1/batch/photos/edit request handler.
		BatchPhotosEdit(router)

		// Specify the unique IDs of the photos used for testing.
		photoUIDs := `["pqkm36fjqvset9uy", "pqkm36fjqvset9uz"]`

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

		// Check the edit response values.
		editPhotos := gjson.Get(editBody, "models").Array()
		assert.Equal(t, len(editPhotos), 2)
		// Send the edit form values back to the same API endpoint and check for errors.
		saveResponse := PerformRequestWithBody(app,
			"POST", "/api/v1/batch/photos/edit",
			fmt.Sprintf(`{"photos": %s, "values": %s}`, photoUIDs,
				"{"+
					"\"Altitude\":{\"value\":0,\"mixed\":false,\"action\":\"update\"},"+
					"\"Year\":{\"value\":-1,\"mixed\":false,\"action\":\"update\"},"+
					"\"Month\":{\"value\":-1,\"mixed\":false,\"action\":\"update\"},"+
					"\"Day\":{\"value\":-1,\"mixed\":false,\"action\":\"update\"},"+
					"\"Title\":{\"value\":\"\",\"mixed\":false,\"action\":\"remove\"},"+
					"\"Caption\":{\"value\":\"\",\"mixed\":false,\"action\":\"remove\"},"+
					"\"DetailsSubject\":{\"value\":\"\",\"mixed\":false,\"action\":\"remove\"},"+
					"\"DetailsArtist\":{\"value\":\"\",\"mixed\":false,\"action\":\"remove\"},"+
					"\"DetailsCopyright\":{\"value\":\"\",\"mixed\":false,\"action\":\"remove\"},"+
					"\"DetailsLicense\":{\"value\":\"\",\"mixed\":false,\"action\":\"remove\"}"+
					"}"),
		)

		// Check the save response status code.
		assert.Equal(t, http.StatusOK, saveResponse.Code)

		// Check the save response body.
		saveBody := saveResponse.Body.String()
		assert.NotEmpty(t, saveBody)

		// Check the save response values.
		saveValues := gjson.Get(saveBody, "values").Raw
		altitudeAfter := gjson.Get(saveValues, "Altitude")
		assert.Equal(t, "{\"value\":0,\"mixed\":false,\"action\":\"none\"}", altitudeAfter.String())
		yearAfter := gjson.Get(saveValues, "Year")
		assert.Equal(t, "{\"value\":-1,\"mixed\":false,\"action\":\"none\"}", yearAfter.String())
		dayAfter := gjson.Get(saveValues, "Day")
		assert.Equal(t, "{\"value\":-1,\"mixed\":false,\"action\":\"none\"}", dayAfter.String())
		monthAfter := gjson.Get(saveValues, "Month")
		assert.Equal(t, "{\"value\":-1,\"mixed\":false,\"action\":\"none\"}", monthAfter.String())
		titleAfter := gjson.Get(saveValues, "Title")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":false,\"action\":\"none\"}", titleAfter.String())
		captionAfter := gjson.Get(saveValues, "Caption")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":false,\"action\":\"none\"}", captionAfter.String())
		subjectAfter := gjson.Get(saveValues, "DetailsSubject")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":false,\"action\":\"none\"}", subjectAfter.String())
		artistAfter := gjson.Get(saveValues, "DetailsArtist")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":false,\"action\":\"none\"}", artistAfter.String())
		copyrightAfter := gjson.Get(saveValues, "DetailsCopyright")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":false,\"action\":\"none\"}", copyrightAfter.String())
		licenseAfter := gjson.Get(saveValues, "DetailsLicense")
		assert.Equal(t, "{\"value\":\"\",\"mixed\":false,\"action\":\"none\"}", licenseAfter.String())
	})*/
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
