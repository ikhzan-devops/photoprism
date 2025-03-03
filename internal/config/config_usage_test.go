package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/pkg/fs/duf"
)

func TestConfig_Usage(t *testing.T) {
	c := TestConfig()

	c.options.UsageInfo = true
	result := c.Usage()
	assert.GreaterOrEqual(t, result.Used, uint64(60000000))
	t.Logf("Storage Used: %d MB (%d%%), Free: %d MB (%d%%), Total %d MB", result.Used/duf.MB, result.UsedPct, result.Free/duf.MB, result.FreePct, result.Total/duf.MB)

	c.options.UsageInfo = false
	assert.Equal(t, c.Usage().Used, uint64(0))
}

func TestConfig_Quota(t *testing.T) {
	c := TestConfig()

	assert.Equal(t, uint64(0), c.Quota())
	assert.Equal(t, 0, c.QuotaUsers())
}
