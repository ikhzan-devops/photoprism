package photoprism

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/photoprism/photoprism/internal/ai/face"
	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/entity"
)

// TestFaces_Match exercises the end-to-end matching flow with a loaded test configuration.
func TestFaces_Match(t *testing.T) {
	c := config.TestConfig()

	m := NewFaces(c)

	opt := FacesOptions{
		Force:     true,
		Threshold: 1,
	}

	r, err := m.Match(opt)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(r)
}

// TestBuildFaceCandidates validates that we drop non-matchable faces when building the index.
func TestBuildFaceCandidates(t *testing.T) {
	regular := entity.NewFace("", entity.SrcAuto, face.RandomEmbeddings(3, face.RegularFace))
	require.NotNil(t, regular)

	kids := entity.NewFace("", entity.SrcAuto, face.RandomEmbeddings(3, face.KidsFace))
	require.NotNil(t, kids)

	faces := entity.Faces{*regular, *kids}

	index := buildFaceIndex(faces)

	require.Len(t, index.fallback, 1)
	require.Equal(t, regular.ID, index.fallback[0].ref.ID)
}

// TestSelectBestFace ensures the best candidate is returned after indexing.
func TestSelectBestFace(t *testing.T) {
	markerEmb := face.RandomEmbeddings(1, face.RegularFace)

	matchFace := entity.NewFace("", entity.SrcAuto, markerEmb)
	require.NotNil(t, matchFace)

	// Force a different face that should not be a better match.
	otherEmb := face.RandomEmbeddings(4, face.RegularFace)
	otherFace := entity.NewFace("", entity.SrcAuto, otherEmb)
	require.NotNil(t, otherFace)

	faces := entity.Faces{*matchFace, *otherFace}

	index := buildFaceIndex(faces)
	require.Len(t, index.fallback, 2)

	best, dist := selectBestFace(markerEmb, index)
	require.NotNil(t, best)
	require.Equal(t, matchFace.ID, best.ID)
	require.InDelta(t, 0.0, dist, 1e-9)
}
