package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/rnd"
)

// Node represents a registered cluster node persisted to YAML.
type Node struct {
	ID        string            `yaml:"id" json:"id"`
	Name      string            `yaml:"name" json:"name"`
	Type      string            `yaml:"type" json:"type"`
	Labels    map[string]string `yaml:"labels" json:"labels"`
	Internal  string            `yaml:"internalUrl" json:"internalUrl"`
	CreatedAt string            `yaml:"createdAt" json:"createdAt"`
	UpdatedAt string            `yaml:"updatedAt" json:"updatedAt"`
	Secret    string            `yaml:"secret" json:"-"` // never JSON-encoded by default
	SecretRot string            `yaml:"nodeSecretLastRotatedAt" json:"nodeSecretLastRotatedAt"`
	DB        struct {
		Name  string `yaml:"name" json:"name"`
		User  string `yaml:"user" json:"user"`
		RotAt string `yaml:"lastRotatedAt" json:"dbLastRotatedAt"`
	} `yaml:"db" json:"db"`
}

func (n *Node) CloneForResponse() Node {
	cp := *n
	cp.Secret = ""
	return cp
}

type FileRegistry struct {
	conf *config.Config
	dir  string
}

func NewFileRegistry(conf *config.Config) (*FileRegistry, error) {
	dir := filepath.Join(conf.PortalConfigPath(), "nodes")
	if err := fs.MkdirAll(dir); err != nil {
		return nil, err
	}
	return &FileRegistry{conf: conf, dir: dir}, nil
}

func (r *FileRegistry) fileName(id string) string { return filepath.Join(r.dir, id+".yaml") }

func (r *FileRegistry) Put(n *Node) error {
	if n.ID == "" {
		n.ID = rnd.UUID()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if n.CreatedAt == "" {
		n.CreatedAt = now
	}
	n.UpdatedAt = now
	b, err := yaml.Marshal(n)
	if err != nil {
		return err
	}
	return os.WriteFile(r.fileName(n.ID), b, 0o600)
}

func (r *FileRegistry) Get(id string) (*Node, error) {
	b, err := os.ReadFile(r.fileName(id))
	if err != nil {
		return nil, err
	}
	var n Node
	if err = yaml.Unmarshal(b, &n); err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *FileRegistry) FindByName(name string) (*Node, error) {
	entries, err := os.ReadDir(r.dir)
	if err != nil {
		return nil, err
	}
	var best *Node
	var bestTime time.Time
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".yaml" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(r.dir, e.Name()))
		if err != nil || len(b) == 0 {
			continue
		}
		var n Node
		if yaml.Unmarshal(b, &n) == nil && n.Name == name {
			// prefer most recently updated
			if t, _ := time.Parse(time.RFC3339, n.UpdatedAt); best == nil || t.After(bestTime) {
				cp := n
				best = &cp
				bestTime = t
			}
		}
	}
	if best == nil {
		return nil, os.ErrNotExist
	}
	return best, nil
}

// List returns all registered nodes (without filtering), sorted by UpdatedAt descending.
func (r *FileRegistry) List() ([]Node, error) {
	entries, err := os.ReadDir(r.dir)
	if err != nil {
		return nil, err
	}
	result := make([]Node, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".yaml" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(r.dir, e.Name()))
		if err != nil || len(b) == 0 {
			continue
		}
		var n Node
		if yaml.Unmarshal(b, &n) == nil {
			result = append(result, n)
		}
	}
	// Sort by UpdatedAt desc if possible (RFC3339 timestamps or empty)
	sort.Slice(result, func(i, j int) bool {
		ti, _ := time.Parse(time.RFC3339, result[i].UpdatedAt)
		tj, _ := time.Parse(time.RFC3339, result[j].UpdatedAt)
		return ti.After(tj)
	})
	return result, nil
}

// Delete removes a node file by id.
func (r *FileRegistry) Delete(id string) error {
	if id == "" {
		return os.ErrNotExist
	}
	return os.Remove(r.fileName(id))
}

func (r *FileRegistry) RotateSecret(id string) (*Node, error) {
	n, err := r.Get(id)
	if err != nil {
		return nil, err
	}
	n.Secret = rnd.Base62(48)
	n.SecretRot = time.Now().UTC().Format(time.RFC3339)
	if err = r.Put(n); err != nil {
		return nil, err
	}
	return n, nil
}

// MarshalJSON customizes JSON output to include nested db fields inline in some responses if needed.
func (n Node) MarshalJSON() ([]byte, error) {
	type Alias Node
	return json.Marshal(Alias(n))
}
