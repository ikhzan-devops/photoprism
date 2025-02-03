package face

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/photoprism/photoprism/pkg/vector"
)

// Embedding represents a face embedding.
type Embedding struct {
	Vector vector.Vector
}

// Dim defines the number of face embedding vector dimensions.
const Dim = 512

var (
	NullVector    = make(vector.Vector, Dim)
	NullEmbedding = Embedding{Vector: NullVector}
)

// NewEmbedding creates a new embedding from an inference result.
func NewEmbedding(values interface{}) Embedding {
	if values == nil {
		return NullEmbedding
	} else if v, err := vector.NewVector(values); err != nil {
		return NullEmbedding
	} else {
		return Embedding{Vector: v}
	}
}

// Null checks if this is a null embedding.
func (m Embedding) Null() bool {
	return len(m.Vector) == 0
}

// Dim returns the dimensions of the embedded vector.
func (m Embedding) Dim() int {
	return len(m.Vector)
}

// Kind returns the type of face e.g. regular, kids, or ignored.
func (m Embedding) Kind() Kind {
	if m.KidsFace() {
		return KidsFace
	} else if m.Ignored() {
		return IgnoredFace
	}

	return RegularFace
}

// SkipMatching checks if the face embedding seems unsuitable for matching.
func (m Embedding) SkipMatching() bool {
	return m.KidsFace() || m.Ignored()
}

// CanMatch tests if the face embedding is not excluded.
func (m Embedding) CanMatch() bool {
	return !m.Ignored()
}

// Dist calculates the distance to another face embedding.
func (m Embedding) Dist(other Embedding) float64 {
	if len(other.Vector) == 0 || len(m.Vector) != len(other.Vector) {
		return -1
	}

	// TODO: Use CosineDist()
	return m.Vector.EuclideanDist(other.Vector)
}

// Norm returns the face embedding vector size (magnitude),
// see https://builtin.com/data-science/vector-norms.
func (m Embedding) Norm() float64 {
	return m.Vector.EuclideanNorm()
}

// MarshalJSON returns the face embedding as JSON.
func (m Embedding) MarshalJSON() ([]byte, error) {
	if len(m.Vector) < 1 {
		return []byte(""), nil
	}

	if result, err := json.Marshal(m.Vector); err != nil {
		return []byte(""), err
	} else {
		return result, nil
	}
}

// UnmarshalJSON sets the embedding vector as JSON.
func (m Embedding) UnmarshalJSON(b []byte) error {
	if len(b) < 1 {
		return nil
	}

	return json.Unmarshal(b, &m.Vector)
}

// JSON returns the face embedding as JSON bytes.
func (m Embedding) JSON() []byte {
	result, _ := m.MarshalJSON()
	return result
}

// UnmarshalEmbedding parses a single face embedding JSON.
func UnmarshalEmbedding(s string) (result Embedding, err error) {
	if s == "" {
		return result, fmt.Errorf("cannot unmarshal embedding, empty string provided")
	} else if !strings.HasPrefix(s, "[") {
		return result, fmt.Errorf("cannot unmarshal embedding, invalid json provided")
	}

	var v = make([]float64, Dim)

	err = json.Unmarshal([]byte(s), &v)

	if err != nil {
		return NewEmbedding(v), err
	}

	return NewEmbedding(v), nil
}
