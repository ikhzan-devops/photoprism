package vision

import (
	"sync"

	"github.com/photoprism/photoprism/internal/ai/classify"
	"github.com/photoprism/photoprism/pkg/clean"
)

var modelMutex = sync.Mutex{}

// Model represents a computer vision model configuration.
type Model struct {
	Name          string   `yaml:"Name,omitempty" json:"name,omitempty"`
	Version       string   `yaml:"Version,omitempty" json:"version,omitempty"`
	Resolution    int      `yaml:"Resolution,omitempty" json:"resolution,omitempty"`
	Url           string   `yaml:"Url,omitempty" json:"-"`
	Method        string   `yaml:"Method,omitempty" json:"-"`
	Format        string   `yaml:"Format,omitempty" json:"-"`
	Path          string   `yaml:"Path,omitempty" json:"-"`
	Tags          []string `yaml:"Tags,omitempty" json:"-"`
	Disabled      bool     `yaml:"Disabled,omitempty" json:"-"`
	classifyModel *classify.Model
}

// Models represents a set of computer vision models.
type Models []*Model

// ClassifyModel returns the matching classify model instance, if any.
func (m *Model) ClassifyModel() *classify.Model {
	modelMutex.Lock()
	defer modelMutex.Unlock()

	if m.classifyModel != nil {
		return m.classifyModel
	}

	switch m.Name {
	case "":
		log.Warnf("vision: missing name, model instance cannot be created")
		return nil
	case NasnetModel.Name, "nasnet":
		if model := classify.NewNasnet(AssetsPath, m.Disabled); model == nil {
			return nil
		} else if err := model.Init(); err != nil {
			log.Errorf("vision: %s (init nasnet model)", err)
			return nil
		} else {
			m.classifyModel = model
		}
	default:
		if m.Path == "" {
			m.Path = clean.TypeLowerUnderscore(m.Name)
		}

		if m.Resolution <= 0 {
			m.Resolution = DefaultResolution
		}

		if model := classify.NewModel(AssetsPath, m.Path, m.Resolution, m.Tags, m.Disabled); model == nil {
			return nil
		} else if err := model.Init(); err != nil {
			log.Errorf("vision: %s (init %s)", err, m.Path)
			return nil
		} else {
			m.classifyModel = model
		}
	}

	return m.classifyModel
}
