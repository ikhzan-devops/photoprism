package config

import (
	"math"
	"time"

	gc "github.com/patrickmn/go-cache"

	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/fs/duf"
)

var usageInfoCache = gc.New(5*time.Minute, 5*time.Minute)

// FlushUsageInfoCache resets the usage information cache.
func FlushUsageInfoCache() {
	usageInfoCache.Flush()
}

// Usage represents storage usage information.
type Usage struct {
	Used    uint64 `json:"used"`
	UsedPct int    `json:"usedPct"`
	Free    uint64 `json:"free"`
	FreePct int    `json:"freePct"`
	Total   uint64 `json:"total"`
}

// Usage returns the used, free and total storage size in bytes and caches the result.
func (c *Config) Usage() Usage {
	// Return nil if feature is not enabled.
	if !c.UsageInfo() {
		return Usage{}
	}

	originalsPath := c.OriginalsPath()

	if cached, hit := usageInfoCache.Get(originalsPath); hit && cached != nil {
		return cached.(Usage)
	}

	info := Usage{}

	if err := c.Db().Unscoped().
		Table("files").
		Select("SUM(file_size) AS used").
		Where("deleted_at IS NULL").
		Take(&info).Error; err != nil {
		log.Warnf("config: failed to calculate indexed file usage (%s)", err.Error())
	}

	quotaTotal := c.QuotaBytes()

	if m, err := duf.PathInfo(originalsPath); err == nil {
		info.Free = m.Free
		info.Total = info.Used + m.Free
	} else {
		log.Debugf("config: failed to detect filesystem usage (%s)", err.Error())
	}

	if quotaTotal > 0 && quotaTotal < info.Total {
		info.Total = quotaTotal
	}

	if info.Total > 0 {
		info.UsedPct = int(math.RoundToEven(float64(info.Used) / float64(info.Total) * 100))
	}

	if info.Used > 0 && info.UsedPct <= 0 {
		info.UsedPct = 1
	}

	info.FreePct = 100 - info.UsedPct

	usageInfoCache.SetDefault(originalsPath, info)

	return info
}

// UsageInfo returns true if resource usage information should be displayed in the user interface.
func (c *Config) UsageInfo() bool {
	return c.options.UsageInfo || c.options.Quota > 0
}

// Quota returns the maximum aggregated size of all indexed files in megabytes, or 0 if no quota exists.
func (c *Config) Quota() uint64 {
	if c.options.Quota <= 0 {
		return 0
	}

	return c.options.Quota
}

// QuotaBytes returns the maximum aggregated size of all indexed files in bytes, or 0 if no quota exists.
func (c *Config) QuotaBytes() uint64 {
	if c.options.Quota <= 0 {
		return 0
	}

	return c.options.Quota * fs.MB
}

// QuotaUsers returns the maximum number of registered user accounts, or 0 if no quota exists.
func (c *Config) QuotaUsers() int {
	if c.options.QuotaUsers <= 0 {
		return 0
	}

	return c.options.QuotaUsers
}
