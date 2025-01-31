package video

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/header"
)

func TestContentType(t *testing.T) {
	t.Run("QuickTime", func(t *testing.T) {
		assert.Equal(t, fs.MimeTypeMP4, ContentType(fs.MimeTypeMOV, "", ""))
	})
	t.Run("QuickTime_HVC", func(t *testing.T) {
		assert.Equal(t, header.ContentTypeHEVC, ContentType(fs.MimeTypeMOV, "mov", CodecHEVC))
	})
	t.Run("MP4", func(t *testing.T) {
		assert.Equal(t, fs.MimeTypeMP4, ContentType(fs.MimeTypeMP4, "", ""))
	})
	t.Run("MP4_AVC", func(t *testing.T) {
		assert.Equal(t, header.ContentTypeAVC, ContentType(fs.MimeTypeMP4, "", CodecAVC))
	})
	t.Run("MP4_HVC", func(t *testing.T) {
		assert.Equal(t, header.ContentTypeHEVC, ContentType(fs.MimeTypeMP4, "", CodecHEVC))
	})
}
