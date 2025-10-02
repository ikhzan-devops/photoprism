package photoprism

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize/english"

	"github.com/photoprism/photoprism/internal/ai/face"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/entity/query"
)

// FacesMatchResult represents the outcome of Faces.Match().
type FacesMatchResult struct {
	Updated    int64
	Recognized int64
	Unknown    int64
}

type faceCandidate struct {
	ref             *entity.Face
	emb             face.Embedding
	sampleRadius    float64
	collisionRadius float64
}

// Add adds result counts.
func (r *FacesMatchResult) Add(result FacesMatchResult) {
	r.Updated += result.Updated
	r.Recognized += result.Recognized
	r.Unknown += result.Unknown
}

// buildFaceCandidates filters the provided faces down to a slice that can be used for
// distance-based matching while caching the embeddings we repeatedly compare against.
func buildFaceCandidates(faces entity.Faces) []faceCandidate {
	candidates := make([]faceCandidate, 0, len(faces))

	for i := range faces {
		f := &faces[i]

		if f.SkipMatching() {
			continue
		}

		embedding := f.Embedding()

		if len(embedding) == 0 {
			continue
		}

		candidates = append(candidates, faceCandidate{
			ref:             f,
			emb:             embedding,
			sampleRadius:    f.SampleRadius,
			collisionRadius: f.CollisionRadius,
		})
	}

	return candidates
}

// match checks whether the supplied marker embeddings fall within the distance and collision
// thresholds for the candidate face, returning the match flag and distance.
func (c faceCandidate) match(embeddings face.Embeddings) (bool, float64) {
	if embeddings.Empty() || len(c.emb) == 0 {
		return false, -1
	}

	dist := minEmbeddingDistance(c.emb, embeddings)

	if dist < 0 {
		return false, dist
	}

	if dist > (c.sampleRadius + face.MatchDist) {
		return false, dist
	}

	if c.collisionRadius > 0.1 && dist > c.collisionRadius {
		return false, dist
	}

	return true, dist
}

// selectBestFace finds the best matching face candidate for the given marker embeddings.
func selectBestFace(embeddings face.Embeddings, candidates []faceCandidate) (*entity.Face, float64) {
	var best *entity.Face
	bestDist := -1.0

	for i := range candidates {
		if ok, dist := candidates[i].match(embeddings); ok {
			if best == nil || dist < bestDist {
				best = candidates[i].ref
				bestDist = dist
			}
		}
	}

	return best, bestDist
}

// Match matches markers with faces and subjects.
func (w *Faces) Match(opt FacesOptions) (result FacesMatchResult, err error) {
	if w.Disabled() {
		return result, fmt.Errorf("face recognition is disabled")
	}

	var unmatchedMarkers int

	// Skip matching if index contains no new face markers, and force option isn't set.
	if opt.Force {
		log.Infof("faces: updating all markers")
	} else if unmatchedMarkers = query.CountUnmatchedFaceMarkers(); unmatchedMarkers > 0 {
		log.Infof("faces: found %s", english.Plural(unmatchedMarkers, "unmatched marker", "unmatched markers"))
	} else {
		log.Debugf("faces: found no unmatched markers")
	}

	matchedAt := entity.TimeStamp()

	if opt.Force || unmatchedMarkers > 0 {
		faces, err := query.Faces(false, false, false, false)

		if err != nil {
			return result, err
		}

		if r, err := w.MatchFaces(faces, opt.Force, nil); err != nil {
			return result, err
		} else {
			result.Add(r)
		}
	}

	// Find unmatched faces.
	if unmatchedFaces, err := query.Faces(false, true, false, false); err != nil {
		log.Error(err)
	} else if len(unmatchedFaces) > 0 {
		if r, err := w.MatchFaces(unmatchedFaces, false, matchedAt); err != nil {
			return result, err
		} else {
			result.Add(r)
		}

		for _, m := range unmatchedFaces {
			if err := m.Matched(); err != nil {
				log.Warnf("faces: %s (update match timestamp)", err)
			}
		}
	}

	// Update remaining markers based on previous matches.
	if m, err := query.MatchFaceMarkers(); err != nil {
		return result, err
	} else {
		result.Recognized += m
	}

	return result, nil
}

// MatchFaces matches markers against a slice of faces.
func (w *Faces) MatchFaces(faces entity.Faces, force bool, matchedBefore *time.Time) (result FacesMatchResult, err error) {
	limit := 500

	candidates := buildFaceCandidates(faces)

	if len(candidates) == 0 {
		log.Debugf("faces: no eligible faces for matching")
		return result, nil
	}

	max := query.CountMarkers(entity.MarkerFace)
	processed := make(map[string]struct{}, max)
	totalProcessed := 0

	offset := 0

	for {
		var markers entity.Markers

		if force {
			markers, err = query.FaceMarkers(limit, offset)
		} else {
			markers, err = query.UnmatchedFaceMarkers(limit, 0, matchedBefore)
		}

		if err != nil {
			return result, err
		}

		if len(markers) == 0 {
			break
		}

		if force {
			offset += len(markers)
			if offset >= max {
				offset = max
			}
		}

		batchProcessed := 0

		for _, marker := range markers {
			if _, seen := processed[marker.MarkerUID]; seen {
				continue
			}

			processed[marker.MarkerUID] = struct{}{}
			totalProcessed++
			batchProcessed++

			if w.Canceled() {
				return result, fmt.Errorf("worker canceled")
			}

			// Skip invalid markers.
			if marker.MarkerInvalid || marker.MarkerType != entity.MarkerFace || len(marker.EmbeddingsJSON) == 0 {
				continue
			}

			markerEmbeddings := marker.Embeddings()

			if markerEmbeddings.Empty() {
				continue
			}

			// Pointer to the matching face.
			selFace, dist := selectBestFace(markerEmbeddings, candidates)

			// Marker already has the best matching face?
			if !marker.HasFace(selFace, dist) {
				// Marker needs a (new) face.
			} else {
				log.Debugf("faces: marker %s already has the best matching face %s with dist %f", marker.MarkerUID, marker.FaceID, marker.FaceDist)

				if err := marker.Matched(); err != nil {
					log.Warnf("faces: %s while updating marker %s match timestamp", err, marker.MarkerUID)
				}

				continue
			}

			// No matching face?
			if selFace == nil {
				if updated, err := marker.ClearFace(); err != nil {
					log.Warnf("faces: %s (clear marker face)", err)
				} else if updated {
					result.Updated++
				}

				continue
			}

			// Assign matching face to marker.
			updated, err := marker.SetFace(selFace, dist)

			if err != nil {
				log.Warnf("faces: %s while setting a face for marker %s", err, marker.MarkerUID)
				continue
			}

			if updated {
				result.Updated++
			}

			if marker.SubjUID != "" {
				result.Recognized++
			} else {
				result.Unknown++
			}
		}

		if batchProcessed == 0 {
			log.Debugf("faces: no new markers to match, stopping")
			break
		}

		log.Debugf("faces: matched %s", english.Plural(totalProcessed, "marker", "markers"))

		if totalProcessed >= max {
			break
		}

		time.Sleep(50 * time.Millisecond)
	}

	return result, err
}
