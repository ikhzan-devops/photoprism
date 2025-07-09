package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType_IsMain(t *testing.T) {
	t.Run("Unknown", func(t *testing.T) {
		assert.False(t, Unknown.IsMain())
	})
	t.Run("Image", func(t *testing.T) {
		assert.True(t, Image.IsMain())
	})
	t.Run("Video", func(t *testing.T) {
		assert.True(t, Video.IsMain())
	})
	t.Run("Sidecar", func(t *testing.T) {
		assert.False(t, Sidecar.IsMain())
	})
}

func TestType_Priority(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		assert.Equal(t, Priority[Unknown], Priority["foo"])
		assert.Equal(t, Priority[Image], Priority[Image])
		assert.Equal(t, Priority[Live], Priority[Live])
		assert.Equal(t, Priority[Animated], Priority[Animated])
		assert.Equal(t, Priority[Animated], Priority[Audio])
		assert.Equal(t, Priority[Animated], Priority[Document])
	})
	t.Run("Less", func(t *testing.T) {
		assert.Less(t, Priority[Unknown], Priority[Image])
		assert.Less(t, Priority[Unknown], Priority[Live])
		assert.Less(t, Priority[Image], Priority[Live])
		assert.Less(t, Priority[Video], Priority[Live])
	})
}

func TestType_Unknown(t *testing.T) {
	t.Run("Unknown", func(t *testing.T) {
		assert.True(t, Unknown.Unknown())
	})
	t.Run("Image", func(t *testing.T) {
		assert.False(t, Image.Unknown())
	})
	t.Run("Video", func(t *testing.T) {
		assert.False(t, Video.Unknown())
	})
	t.Run("Sidecar", func(t *testing.T) {
		assert.False(t, Sidecar.Unknown())
	})
}

func TestType_Equal(t *testing.T) {
	t.Run("UnknownUnknown", func(t *testing.T) {
		assert.True(t, Unknown.Equal(""))
	})
	t.Run("ImageImage", func(t *testing.T) {
		assert.True(t, Image.Equal(Image.String()))
	})
	t.Run("VideoImage", func(t *testing.T) {
		assert.False(t, Video.Equal(Image.String()))
	})
	t.Run("SidecarUnknown", func(t *testing.T) {
		assert.False(t, Sidecar.Equal(Unknown.String()))
	})
}

func TestType_NotEqual(t *testing.T) {
	t.Run("UnknownUnknown", func(t *testing.T) {
		assert.False(t, Unknown.NotEqual(""))
	})
	t.Run("ImageImage", func(t *testing.T) {
		assert.False(t, Image.NotEqual(Image.String()))
	})
	t.Run("VideoImage", func(t *testing.T) {
		assert.True(t, Video.NotEqual(Image.String()))
	})
	t.Run("SidecarUnknown", func(t *testing.T) {
		assert.True(t, Sidecar.NotEqual(Unknown.String()))
	})
}
