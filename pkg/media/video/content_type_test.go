package video

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/pkg/media/http/header"
)

func TestContentType(t *testing.T) {
	t.Run("QuickTime", func(t *testing.T) {
		assert.Equal(t, header.ContentTypeMp4, ContentType(header.ContentTypeMov, "", ""))
	})
	t.Run("QuickTime/Hvc", func(t *testing.T) {
		assert.Equal(t, header.ContentTypeMp4HvcMain, ContentType(header.ContentTypeMov, "mov", CodecHvc))
	})
	t.Run("Mp4", func(t *testing.T) {
		assert.Equal(t, header.ContentTypeMp4, ContentType(header.ContentTypeMp4, "", ""))
	})
	t.Run("Mp4/Avc", func(t *testing.T) {
		assert.Equal(t, header.ContentTypeMp4AvcMain, ContentType(header.ContentTypeMp4, "", CodecAvc))
	})
	t.Run("Mp4/Hvc", func(t *testing.T) {
		assert.Equal(t, header.ContentTypeMp4HvcMain, ContentType(header.ContentTypeMp4, "", CodecHvc))
	})
}

func TestCompatible(t *testing.T) {
	t.Run("True", func(t *testing.T) {
		assert.True(t, Compatible(header.ContentTypeJpeg, header.ContentTypeJpeg))
		assert.True(t, Compatible(header.ContentTypeMp4, header.ContentTypeMov))
		assert.True(t, Compatible(header.ContentTypeMp4Hvc, header.ContentTypeMp4HvcMain))
		assert.True(t, Compatible(header.ContentTypeMp4Hvc, header.ContentTypeMp4HvcMain10))
		assert.True(t, Compatible(header.ContentTypeMp4Avc, header.ContentTypeMp4AvcHigh))
	})
	t.Run("False", func(t *testing.T) {
		assert.False(t, Compatible("", ""))
		assert.False(t, Compatible("", header.ContentTypeMov))
		assert.False(t, Compatible(header.ContentTypeMp4, header.ContentTypeMp4Avc))
	})
}
