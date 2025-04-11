package vision

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/photoprism/photoprism/internal/ai/classify"
	"github.com/photoprism/photoprism/internal/ai/face"
	"github.com/photoprism/photoprism/internal/ai/nsfw"
	"github.com/photoprism/photoprism/internal/ai/tensorflow"
	"github.com/photoprism/photoprism/pkg/clean"
)

var modelMutex = sync.Mutex{}

// Model represents a computer vision model configuration.
type Model struct {
	Type          ModelType             `yaml:"Type,omitempty" json:"type,omitempty"`
	Name          string                `yaml:"Name,omitempty" json:"name,omitempty"`
	Version       string                `yaml:"Version,omitempty" json:"version,omitempty"`
	Resolution    int                   `yaml:"Resolution,omitempty" json:"resolution,omitempty"`
	Meta          *tensorflow.ModelInfo `yaml:"Meta,omitempty" json:"meta,omitempty"`
	Uri           string                `yaml:"Uri,omitempty" json:"-"`
	Key           string                `yaml:"Key,omitempty" json:"-"`
	Method        string                `yaml:"Method,omitempty" json:"-"`
	Path          string                `yaml:"Path,omitempty" json:"-"`
	Disabled      bool                  `yaml:"Disabled,omitempty" json:"-"`
	classifyModel *classify.Model
	faceModel     *face.Model
	nsfwModel     *nsfw.Model
}

// Models represents a set of computer vision models.
type Models []*Model

// Endpoint returns the remote service request method and endpoint URL, if any.
func (m *Model) Endpoint() (uri, method string) {
	if m.Uri == "" && ServiceUri == "" || m.Type == "" {
		return "", ""
	}

	if m.Method != "" {
		method = m.Method
	} else {
		method = http.MethodPost
	}

	if m.Uri != "" {
		return m.Uri, method
	} else {
		return fmt.Sprintf("%s/%s", ServiceUri, clean.TypeLowerUnderscore(m.Type)), method
	}
}

// EndpointKey returns the access token belonging to the remote service endpoint, if any.
func (m *Model) EndpointKey() string {
	if m.Key != "" {
		return m.Key
	} else if ServiceKey != "" {
		return ServiceKey
	}

	return ""
}

// ClassifyModel returns the matching classify model instance, if any.
func (m *Model) ClassifyModel() *classify.Model {
	// Use mutex to prevent models from being loaded and
	// initialized twice by different indexing workers.
	modelMutex.Lock()
	defer modelMutex.Unlock()

	// Return the existing model instance if it has already been created.
	if m.classifyModel != nil {
		return m.classifyModel
	}

	switch m.Name {
	case "":
		log.Warnf("vision: missing name, model instance cannot be created")
		return nil
	case NasnetModel.Name, "nasnet":
		// Load and initialize the Nasnet image classification model.
		if model := classify.NewNasnet(AssetsPath, m.Disabled); model == nil {
			return nil
		} else if err := model.Init(); err != nil {
			log.Errorf("vision: %s (init nasnet model)", err)
			return nil
		} else {
			m.classifyModel = model
		}
	default:
		// Set model path from model name if no path is configured.
		if m.Path == "" {
			m.Path = clean.TypeLowerUnderscore(m.Name)
		}

		if m.Meta == nil {
			m.Meta = &tensorflow.ModelInfo{}
		}

		// Set default thumbnail resolution if no tags are configured.
		if m.Resolution <= 0 {
			m.Resolution = DefaultResolution
		} else {
			if m.Meta.Input == nil {
				m.Meta.Input = new(tensorflow.PhotoInput)
			}

			m.Meta.Input.SetResolution(m.Resolution)
			m.Meta.Input.Channels = 3
		}

		// Try to load custom model based on the configuration values.
		if model := classify.NewModel(AssetsPath, m.Path, m.Meta, m.Disabled); model == nil {
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

// FaceModel returns the matching face model instance, if any.
func (m *Model) FaceModel() *face.Model {
	// Use mutex to prevent models from being loaded and
	// initialized twice by different indexing workers.
	modelMutex.Lock()
	defer modelMutex.Unlock()

	// Return the existing model instance if it has already been created.
	if m.faceModel != nil {
		return m.faceModel
	}

	switch m.Name {
	case "":
		log.Warnf("vision: missing name, model instance cannot be created")
		return nil
	case FacenetModel.Name, "facenet":
		// Load and initialize the Nasnet image classification model.
		if model := face.NewModel(FaceNetModelPath, CachePath, m.Resolution, m.Meta.Tags, m.Disabled); model == nil {
			return nil
		} else if err := model.Init(); err != nil {
			log.Errorf("vision: %s (init %s)", err, m.Path)
			return nil
		} else {
			m.faceModel = model
		}
	default:
		// Set model path from model name if no path is configured.
		if m.Path == "" {
			m.Path = clean.TypeLowerUnderscore(m.Name)
		}

		// Set default thumbnail resolution if no tags are configured.
		if m.Resolution <= 0 {
			m.Resolution = DefaultResolution
		}

		if m.Meta == nil {
			m.Meta = &tensorflow.ModelInfo{}
		}

		// Set default tag if no tags are configured.
		if len(m.Meta.Tags) == 0 {
			m.Meta.Tags = []string{"serve"}
		}

		// Try to load custom model based on the configuration values.
		if model := face.NewModel(filepath.Join(AssetsPath, m.Path), CachePath, m.Resolution, m.Meta.Tags, m.Disabled); model == nil {
			return nil
		} else if err := model.Init(); err != nil {
			log.Errorf("vision: %s (init %s)", err, m.Path)
			return nil
		} else {
			m.faceModel = model
		}
	}

	return m.faceModel
}

// NsfwModel returns the matching nsfw model instance, if any.
func (m *Model) NsfwModel() *nsfw.Model {
	// Use mutex to prevent models from being loaded and
	// initialized twice by different indexing workers.
	modelMutex.Lock()
	defer modelMutex.Unlock()

	// Return the existing model instance if it has already been created.
	if m.nsfwModel != nil {
		return m.nsfwModel
	}

	switch m.Name {
	case "":
		log.Warnf("vision: missing name, model instance cannot be created")
		return nil
	case NsfwModel.Name, "nsfw":
		// Load and initialize the Nasnet image classification model.
		if model := nsfw.NewModel(NsfwModelPath, NsfwModel.Meta, m.Disabled); model == nil {
			return nil
		} else if err := model.Init(); err != nil {
			log.Errorf("vision: %s (init %s)", err, m.Path)
			return nil
		} else {
			m.nsfwModel = model
		}
	default:
		// Set model path from model name if no path is configured.
		if m.Path == "" {
			m.Path = clean.TypeLowerUnderscore(m.Name)
		}

		// Set default thumbnail resolution if no tags are configured.
		if m.Resolution <= 0 {
			m.Resolution = DefaultResolution
		} else {
			if m.Meta.Input == nil {
				m.Meta.Input = new(tensorflow.PhotoInput)
			}

			m.Meta.Input.SetResolution(m.Resolution)
			m.Meta.Input.Channels = 3
		}

		if m.Meta == nil {
			m.Meta = &tensorflow.ModelInfo{}
		}

		// Set default tag if no tags are configured.
		if len(m.Meta.Tags) == 0 {
			m.Meta.Tags = []string{"serve"}
		}

		// Try to load custom model based on the configuration values.
		if model := nsfw.NewModel(filepath.Join(AssetsPath, m.Path), m.Meta, m.Disabled); model == nil {
			return nil
		} else if err := model.Init(); err != nil {
			log.Errorf("vision: %s (init %s)", err, m.Path)
			return nil
		} else {
			m.nsfwModel = model
		}
	}

	return m.nsfwModel
}
