package face

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/photoprism/photoprism/pkg/vector"
)

// Embeddings represents a face embedding cluster.
type Embeddings []Embedding

// NewEmbeddings creates a new embeddings from float64 slices.
func NewEmbeddings(values [][]float64) Embeddings {
	result := make(Embeddings, 0, len(values))

	var i int

	for i = range values {
		result = append(result, NewEmbedding(values[i]))
	}

	return result
}

// NewEmbeddingsFromInference creates a new embeddings from float32 inference slices.
func NewEmbeddingsFromInference(values [][]float32) Embeddings {
	result := make(Embeddings, len(values))

	var v []float32
	var i int

	for i, v = range values {
		e := NewEmbedding(v)

		if e.CanMatch() {
			result[i] = e
		}
	}

	return result
}

// Empty tests if embeddings are empty.
func (embeddings Embeddings) Empty() bool {
	if len(embeddings) < 1 {
		return true
	}

	return embeddings[0].Dim() < 1
}

// Count returns the number of embeddings.
func (embeddings Embeddings) Count() int {
	if embeddings.Empty() {
		return 0
	}

	return len(embeddings)
}

// Kind returns the type of face e.g. regular, kids, or ignored.
func (embeddings Embeddings) Kind() (result Kind) {
	for _, e := range embeddings {
		if k := e.Kind(); k > result {
			result = k
		}
	}

	return result
}

// One tests if there is exactly one embedding.
func (embeddings Embeddings) One() bool {
	return embeddings.Count() == 1
}

// First returns the first face embedding.
func (embeddings Embeddings) First() Embedding {
	if embeddings.Empty() {
		return NullEmbedding
	}

	return embeddings[0]
}

// Float64 returns embeddings as a float64 slice.
func (embeddings Embeddings) Float64() [][]float64 {
	result := make([][]float64, len(embeddings))

	for i, e := range embeddings {
		result[i] = e.Vector
	}

	return result
}

// Contains tests if another embeddings is contained within a radius.
func (embeddings Embeddings) Contains(other Embedding, radius float64) bool {
	for _, e := range embeddings {
		if d := e.Dist(other); d < radius {
			return true
		}
	}

	return false
}

// Dist returns the minimum distance to an embedding.
func (embeddings Embeddings) Dist(other Embedding) (dist float64) {
	dist = -1

	for _, e := range embeddings {
		if d := e.Dist(other); d < dist || dist < 0 {
			dist = d
		}
	}

	return dist
}

// JSON returns the embeddings as JSON bytes.
func (embeddings Embeddings) JSON() []byte {
	var noResult = []byte("")

	if embeddings.Empty() {
		return noResult
	}

	if result, err := json.Marshal(embeddings); err != nil {
		return noResult
	} else {
		return result
	}
}

// MarshalJSON returns the face embeddings as JSON.
func (embeddings Embeddings) MarshalJSON() ([]byte, error) {
	values := make(vector.Vectors, len(embeddings))

	for i := range embeddings {
		values[i] = embeddings[i].Vector
	}

	return json.Marshal(values)
}

// EmbeddingsMidpoint returns the embeddings vector midpoint.
func EmbeddingsMidpoint(embeddings Embeddings) (result Embedding, radius float64, count int) {
	// Return if there are no embeddings.
	if embeddings.Empty() {
		return Embedding{}, 0, 0
	}

	// Count embeddings.
	count = len(embeddings)

	// Only one embedding?
	if count == 1 {
		// Return embedding if there is only one.
		return embeddings[0], 0.0, 1
	}

	dim := embeddings[0].Dim()

	// No embedding values?
	if dim == 0 {
		return Embedding{}, 0.0, count
	}

	// Create a new embedding with the given vector dimension.
	result = NewEmbedding(vector.NullVector(dim))

	// Calculate mean values.
	// TODO: Improve to get better matching results.
	for i := 0; i < dim; i++ {
		values := make(vector.Vector, count)

		for j := 0; j < count; j++ {
			values[j] = embeddings[j].Vector[i]
		}

		result.Vector[i] = values.Mean()
	}

	// Radius is the max embedding distance + 0.01 from result.
	for _, emb := range embeddings {
		if d := result.Dist(emb); d > radius {
			radius = d + 0.01
		}
	}

	return result, radius, count
}

// UnmarshalEmbeddings parses face embedding JSON.
func UnmarshalEmbeddings(s string) (result Embeddings, err error) {
	if s == "" {
		return result, fmt.Errorf("cannot unmarshal empeddings, empty string provided")
	} else if !strings.HasPrefix(s, "[[") {
		return result, fmt.Errorf("cannot unmarshal empeddings, invalid json provided")
	}

	var values [][]float64

	if err = json.Unmarshal([]byte(s), &values); err != nil {
		return result, err
	}

	return NewEmbeddings(values), nil
}
