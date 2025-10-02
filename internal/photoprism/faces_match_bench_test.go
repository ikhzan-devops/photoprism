package photoprism

import (
	"testing"

	"github.com/photoprism/photoprism/internal/ai/face"
	"github.com/photoprism/photoprism/internal/entity"
)

func BenchmarkSelectBestFace(b *testing.B) {
	const candidateCount = 32

	faces := make(entity.Faces, 0, candidateCount)

	for i := 0; i < candidateCount; i++ {
		f := entity.NewFace("", entity.SrcAuto, face.RandomEmbeddings(5, face.RegularFace))
		faces = append(faces, *f)
	}

	candidates := buildFaceCandidates(faces)
	markerEmb := face.RandomEmbeddings(3, face.RegularFace)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		selectBestFace(markerEmb, candidates)
	}
}

func BenchmarkSelectBestFaceLegacy(b *testing.B) {
	const candidateCount = 32

	faces := make(entity.Faces, 0, candidateCount)

	for i := 0; i < candidateCount; i++ {
		f := entity.NewFace("", entity.SrcAuto, face.RandomEmbeddings(5, face.RegularFace))
		faces = append(faces, *f)
	}

	markerEmb := face.RandomEmbeddings(3, face.RegularFace)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		legacySelectBestFace(markerEmb, faces)
	}
}

func legacySelectBestFace(marker face.Embeddings, faces entity.Faces) (*entity.Face, float64) {
	var best *entity.Face
	bestDist := -1.0

	for i := range faces {
		if ok, dist := faces[i].Match(marker); ok {
			if best == nil || dist < bestDist {
				best = &faces[i]
				bestDist = dist
			}
		}
	}

	return best, bestDist
}
