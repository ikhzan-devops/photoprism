package photoprism

import (
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/photoprism/photoprism/internal/entity/query"
)

// Files represents a list of already indexed file names and their unix modification timestamps.
type Files struct {
	count int
	files query.FileMap
	mutex sync.RWMutex
}

// NewFiles returns a new Files instance.
func NewFiles() *Files {
	m := &Files{
		files: make(query.FileMap),
	}

	return m
}

// Init lazily loads the indexed file map from the database and stores the initial count.
func (m *Files) Init() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.files) > 0 {
		m.count = len(m.files)
		return nil
	}

	if err := query.PurgeOrphanDuplicates(); err != nil {
		return fmt.Errorf("%s (purge duplicates)", err.Error())
	}

	files, err := query.IndexedFiles()

	if err != nil {
		return fmt.Errorf("%s (find indexed files)", err.Error())
	} else {
		m.files = files
		m.count = len(files)
		return nil
	}
}

// Done clears the in-memory cache so the next index pass reloads a fresh snapshot.
func (m *Files) Done() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if (len(m.files) - m.count) == 0 {
		return
	}

	m.count = 0
	m.files = make(query.FileMap)
}

// Remove evicts a file entry from the cache (e.g. after deletion or re-import).
func (m *Files) Remove(fileName, fileRoot string) {
	key := path.Join(fileRoot, fileName)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.files, key)
}

// Ignore determines whether a file should be skipped based on modification timestamp, updating the cache when needed.
func (m *Files) Ignore(fileName, fileRoot string, modTime time.Time, rescan bool) bool {
	timestamp := modTime.UTC().Truncate(time.Second).Unix()
	key := path.Join(fileRoot, fileName)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if rescan {
		m.files[key] = timestamp
		return false
	}

	mod, ok := m.files[key]

	if ok && mod == timestamp {
		return true
	} else {
		m.files[key] = timestamp
		return false
	}
}

// Indexed checks if a file has already been indexed without mutating the cache.
func (m *Files) Indexed(fileName, fileRoot string, modTime time.Time, rescan bool) bool {
	if rescan {
		return false
	}

	timestamp := modTime.Unix()
	key := path.Join(fileRoot, fileName)

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	mod, ok := m.files[key]

	if ok && mod == timestamp {
		return true
	} else {
		return false
	}
}

// Exists reports whether the given file key is present in the cache.
func (m *Files) Exists(fileName, fileRoot string) bool {
	key := path.Join(fileRoot, fileName)

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if _, ok := m.files[key]; ok {
		return true
	} else {
		return false
	}
}
