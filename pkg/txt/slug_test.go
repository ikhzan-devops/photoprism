package txt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlug(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t, "", Slug(""))
	})
	t.Run("Gates", func(t *testing.T) {
		assert.Equal(t, "william-henry-gates-iii", Slug("William  Henry Gates III"))
	})
	t.Run("Quotes", func(t *testing.T) {
		assert.Equal(t, "william-henry-gates", Slug("william \"HenRy\" gates' "))
	})
	t.Run("Chinese", func(t *testing.T) {
		assert.Equal(t, "chen-zhao", Slug(" 陈  赵"))
	})
	t.Run("Emoji", func(t *testing.T) {
		assert.Equal(t, "_5cpzfdq", Slug("💎"))
		assert.Equal(t, "_5cpzfea", Slug("💐"))
		assert.Equal(t, "_5cpzfea", Slug("   💐   "))
		assert.Equal(t, "_5cpzfdxqt5jja", Slug("💎💐"))
		assert.Equal(t, "photoprism", Slug("PhotoPrism 💎"))
	})
}

func TestSlugToTitle(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t, "", SlugToTitle(""))
	})
	t.Run("Kitten", func(t *testing.T) {
		assert.Equal(t, "Cute-Kitten", SlugToTitle("cute-kitten"))
	})
	t.Run("Emoji", func(t *testing.T) {
		assert.Equal(t, "💎", SlugToTitle("_5cpzfdq"))
		assert.Equal(t, "💐", SlugToTitle("_5cpzfea"))
		assert.Equal(t, "💎💐", SlugToTitle("_5cpzfdxqt5jja"))
		assert.Equal(t, "PhotoPrism", SlugToTitle("photoprism"))
	})
}
