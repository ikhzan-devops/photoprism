package batch

import (
	"fmt"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/entity/query"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/rnd"
)

type Action = string

const (
	ActionNone   Action = "none"
	ActionUpdate Action = "update"
	ActionAdd    Action = "add"
	ActionRemove Action = "remove"
)

// ApplyAlbums adds/removes the given photo to/from albums according to items action.
func ApplyAlbums(photoUID string, albums Items) error {
	var addTargets []string

	for _, it := range albums.Items {
		switch it.Action {
		case ActionAdd:
			// Add by UID if provided, otherwise use title to create/find
			if it.Value != "" {
				addTargets = append(addTargets, it.Value)
			} else if it.Title != "" {
				addTargets = append(addTargets, it.Title)
			}
		case ActionRemove:
			// Remove only if we have a valid album UID
			if rnd.IsUID(it.Value, entity.AlbumUID) {
				if a, err := query.AlbumByUID(it.Value); err != nil {
					log.Debugf("batch: album %s not found for removal: %s", it.Value, err)
				} else if a.HasID() {
					a.RemovePhotos([]string{photoUID})
				}
			}
		}
	}

	if len(addTargets) > 0 {
		if err := entity.AddPhotoToAlbums(photoUID, addTargets); err != nil {
			return err
		}
	}

	return nil
}

// ApplyLabels adds/removes labels on the given photo according to items action.
func ApplyLabels(photo *entity.Photo, labels Items) error {
	if photo == nil || !photo.HasID() {
		return fmt.Errorf("invalid photo")
	}

	// Track if we changed anything to call SaveLabels once
	changed := false

	for _, it := range labels.Items {
		switch it.Action {
		case ActionAdd:
			// Try by UID first
			var labelEntity *entity.Label
			var err error
			if it.Value != "" {
				labelEntity, err = query.LabelByUID(it.Value)
				if err != nil {
					labelEntity = nil
				}
			}
			if labelEntity == nil && it.Title != "" {
				// Create or find by title
				labelEntity = entity.FirstOrCreateLabel(entity.NewLabel(it.Title, 0))
			}

			if labelEntity == nil {
				log.Debugf("batch: could not resolve label to add: value=%s title=%s", it.Value, clean.Log(it.Title))
				continue
			}

			if err := labelEntity.Restore(); err != nil {
				log.Debugf("batch: could not restore label %s: %s", labelEntity.LabelName, err)
			}

			// Ensure 100% confidence (uncertainty 0) and source 'batch'
			if pl := entity.FirstOrCreatePhotoLabel(entity.NewPhotoLabel(photo.ID, labelEntity.ID, 0, entity.SrcBatch)); pl == nil {
				log.Errorf("batch: failed creating photo-label for photo %d and label %d", photo.ID, labelEntity.ID)
			} else {
				// If it already existed with different values, update it
				if pl.Uncertainty != 0 || pl.LabelSrc != entity.SrcBatch {
					pl.Uncertainty = 0
					pl.LabelSrc = entity.SrcBatch
					if err := entity.Db().Save(pl).Error; err != nil {
						log.Errorf("batch: update label to 100%% confidence failed: %s", err)
					} else {
						changed = true
					}
				} else {
					changed = true
				}
			}

		case ActionRemove:
			if it.Value == "" {
				log.Debugf("batch: label remove skipped (uid required): photo=%s title=%s", photo.PhotoUID, clean.Log(it.Title))
				continue
			}

			labelEntity, err := query.LabelByUID(it.Value)
			if err != nil || labelEntity == nil || !labelEntity.HasID() {
				log.Debugf("batch: label not found for removal by uid: photo=%s uid=%s", photo.PhotoUID, it.Value)
				continue
			}

			if pl, err := query.PhotoLabel(photo.ID, labelEntity.ID); err != nil {
				log.Debugf("batch: photo-label not found for removal: photo=%s label_id=%d", photo.PhotoUID, labelEntity.ID)
			} else if pl != nil {
				// Block label from being auto re-added by setting uncertainty to 100 and marking source as 'batch'.
				pl.Uncertainty = 100
				pl.LabelSrc = entity.SrcBatch
				if err := entity.Db().Save(pl).Error; err != nil {
					log.Errorf("batch: block label failed: %s", err)
				} else {
					log.Debugf("batch: blocked label: photo=%s label_id=%d", photo.PhotoUID, labelEntity.ID)
					changed = true
				}
				_ = photo.RemoveKeyword(labelEntity.LabelName)
			}
		}
	}

	if changed {
		// Reload photo to ensure in-memory labels reflect DB changes before saving derived fields
		if reloaded, err := query.PhotoPreloadByUID(photo.PhotoUID); err == nil && reloaded.HasID() {
			if err := (&reloaded).SaveLabels(); err != nil {
				return err
			}
		} else {
			if err := photo.SaveLabels(); err != nil {
				return err
			}
		}
	}

	return nil
}
