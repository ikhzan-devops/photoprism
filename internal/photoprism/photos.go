package photoprism

import (
	"fmt"
	"sync"
	"time"

	"github.com/photoprism/photoprism/internal/entity"

	"github.com/photoprism/photoprism/internal/entity/query"
)

// Photos caches a lookup table from capture timestamp + S2 cell IDs to photo IDs so workers can skip redundant database scans.
type Photos struct {
	count  int
	photos query.PhotoMap
	mutex  sync.RWMutex
}

// NewPhotos constructs an empty Photos cache. Call Init before using Find.
func NewPhotos() *Photos {
	m := &Photos{
		photos: make(query.PhotoMap),
	}

	return m
}

// Init hydrates the cache from the database if it has not been loaded yet.
// Subsequent calls are no-ops unless the internal map was reset.
func (m *Photos) Init() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.photos) > 0 {
		m.count = len(m.photos)
		return nil
	}

	photos, err := query.IndexedPhotos()

	if err != nil {
		return fmt.Errorf("%s (find indexed photos)", err.Error())
	} else {
		m.photos = photos
		m.count = len(photos)
		return nil
	}
}

// Remove evicts a photo from the lookup table when media has been deleted or re-indexed.
func (m *Photos) Remove(takenAt time.Time, cellId string) {
	key := entity.MapKey(takenAt, cellId)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.photos, key)
}

// Find returns the cached photo ID for the given capture time and cell. Zero means no entry exists.
func (m *Photos) Find(takenAt time.Time, cellId string) uint {
	key := entity.MapKey(takenAt, cellId)

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.photos[key]
}
