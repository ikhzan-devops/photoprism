package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFaces_Embeddings(t *testing.T) {
	m := FaceFixtures.Get("joe-biden")
	m1 := FaceFixtures.Get("jane-doe")
	r := Faces{m, m1}.Embeddings()
	len1 := m.Embedding().Dim()
	len2 := m1.Embedding().Dim()
	assert.Equal(t, len1+len2, len(r[0].Vector)+len(r[1].Vector))
}

func TestFaces_IDs(t *testing.T) {
	m := FaceFixtures.Get("joe-biden")
	m1 := FaceFixtures.Get("jane-doe")
	r := Faces{m, m1}.IDs()
	assert.Equal(t, []string{"VF7ANLDET2BKZNT4VQWJMMC6HBEFDOG6", "VF7ANLDET2BKZNT4VQWJMMC6HBEFDOG7"}, r)
}

func TestDeleteOrphanFaces(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		if count, err := DeleteOrphanFaces(); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("deleted %d faces", count)
		}
	})
}
