package pwa

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/pkg/fs"
)

func TestNewIcons(t *testing.T) {
	t.Run("Standard", func(t *testing.T) {
		c := Config{StaticUri: "https://demo-cdn.photoprism.app/static", Icon: "test"}
		result := NewIcons(c)
		assert.NotEmpty(t, result)
		assert.Equal(t, "https://demo-cdn.photoprism.app/static/icons/test/16.png", result[0].Src)
		assert.Equal(t, "image/png", result[0].Type)
		assert.Equal(t, "16x16", result[0].Sizes)
	})
	t.Run("Custom", func(t *testing.T) {
		c := Config{StaticUri: "https://demo-cdn.photoprism.app/static", Icon: "/test.png"}
		result := NewIcons(c)
		assert.NotEmpty(t, result)
		assert.Equal(t, "/test.png", result[0].Src)
		assert.Equal(t, "image/png", result[0].Type)
		assert.Equal(t, "", result[0].Sizes)
	})
	t.Run("Theme", func(t *testing.T) {
		c := Config{StaticUri: "https://demo-cdn.photoprism.app/static", Icon: "/_theme/example.png", ThemePath: fs.Abs("./testdata"), ThemeUri: "/_theme"}
		result := NewIcons(c)
		assert.NotEmpty(t, result)
		assert.Equal(t, "/_theme/example.png", result[0].Src)
		assert.Equal(t, "image/png", result[0].Type)
		assert.Equal(t, "100x67", result[0].Sizes)
	})
}
