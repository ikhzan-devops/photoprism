package vision

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/ai/classify"
	"github.com/photoprism/photoprism/pkg/fs"
)

func TestLabels(t *testing.T) {
	var assetsPath = fs.Abs("../../../assets")
	var examplesPath = assetsPath + "/examples"

	t.Run("Success", func(t *testing.T) {
		result, err := Labels([]string{examplesPath + "/chameleon_lime.jpg"})

		assert.NoError(t, err)
		assert.IsType(t, classify.Labels{}, result)
		assert.Equal(t, 1, len(result))

		t.Log(result)

		assert.Equal(t, "chameleon", result[0].Name)
		assert.Equal(t, 7, result[0].Uncertainty)
	})
	t.Run("InvalidFile", func(t *testing.T) {
		_, err := Labels([]string{examplesPath + "/notexisting.jpg"})
		assert.Error(t, err)
	})
}
