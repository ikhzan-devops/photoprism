package hooks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/event"
)

func TestPayload_Json(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		timeStamp, timeErr := time.Parse(time.RFC3339Nano, "2025-04-11T11:48:58.540199797Z")

		if timeErr != nil {
			t.Fatal(timeErr)
		}

		id := "49b8a329-5aa6-4b76-ba62-bb3adb001817"

		payload := &Payload{
			Type:      "foo.bar",
			Timestamp: timeStamp.UTC(),
			Data: event.Data{
				"id":     id,
				"hello":  "World!",
				"number": 42,
			},
		}

		result := payload.JSON()
		expected := `{"type":"foo.bar","timestamp":"2025-04-11T11:48:58.540199797Z","data":{"hello":"World!","id":"49b8a329-5aa6-4b76-ba62-bb3adb001817","number":42}}`

		assert.Equal(t, expected, string(result))
	})
}
