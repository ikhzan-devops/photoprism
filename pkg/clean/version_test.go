package clean

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Run("MariaDB", func(t *testing.T) {
		assert.Equal(t, "", Version(""))
		assert.Equal(t, "v10.5.0", Version("v10.5.0"))
		assert.Equal(t, "v10.5.1", Version("10.5.1"))
		assert.Equal(t, "v10.5.5", Version("MariaDB-1:10.5.5+maria~focal"))
		assert.Equal(t, "v10.5.7", Version("v10.5.7"))
		assert.Equal(t, "v10.11.10", Version("10.11.10-MariaDB-ubu2204"))
		assert.Equal(t, "v10.11.11", Version("10.11.11-MariaDB-0ubuntu0.24.04.2-log"))
		assert.Equal(t, "v11.4.5", Version("11.4.5-MariaDB-ubu2404"))
		assert.Equal(t, "", Version("11.7"))
		assert.Equal(t, "v11.7.0", Version("11.7.0"))
		assert.Equal(t, "v11.7.1", Version("11.7.1"))
		assert.Equal(t, "v11.7.01", Version("11.7.01"))
		assert.Equal(t, "v24.04.2", Version("24.04.2-MariaDB-0ubuntu0.10.11.11-log"))
	})
	t.Run("Empty", func(t *testing.T) {
		result := Version("")
		assert.Equal(t, "", result)
	})
	t.Run("Invalid", func(t *testing.T) {
		result := Version("45345.356q636")
		assert.Equal(t, "", result)
	})
}
