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
	// Validate photo UID
	if !rnd.IsUID(photoUID, entity.PhotoUID) {
		return fmt.Errorf("invalid photo uid: %s", photoUID)
	}

	var addTargets []string

	for _, it := range albums.Items {
		switch it.Action {
		case ActionAdd:
			// Validate that we have either value or title
			if it.Value == "" && it.Title == "" {
				return fmt.Errorf("album value or title required for add action")
			}

			// Add by UID if provided, otherwise use title to create/find
			if it.Value != "" {
				// If value is provided, validate it's a proper UID format
				if !rnd.IsUID(it.Value, entity.AlbumUID) {
					return fmt.Errorf("invalid album uid format: %s", it.Value)
				}

				// Check if album exists when adding by UID
				if _, err := query.AlbumByUID(it.Value); err != nil {
					return fmt.Errorf("album not found: %s", it.Value)
				}

				addTargets = append(addTargets, it.Value)
			} else if it.Title != "" {
				addTargets = append(addTargets, it.Title)
			}
		case ActionRemove:
			// Validate that we have a value for removal
			if it.Value == "" {
				return fmt.Errorf("album uid required for remove action")
			}

			// Remove only if we have a valid album UID
			if !rnd.IsUID(it.Value, entity.AlbumUID) {
				return fmt.Errorf("invalid album uid format: %s", it.Value)
			}

			if a, err := query.AlbumByUID(it.Value); err != nil {
				return fmt.Errorf("album not found for removal: %s", it.Value)
			} else if a.HasID() {
				a.RemovePhotos([]string{photoUID})
			}
		case ActionNone, ActionUpdate:
			// Valid actions that do nothing for albums
			continue
		default:
			return fmt.Errorf("invalid action: %s", it.Action)
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
			// Validate that we have either value or title
			if it.Value == "" && it.Title == "" {
				return fmt.Errorf("label value or title required for add action")
			}

			// Try by UID first
			var labelEntity *entity.Label
			var err error
			if it.Value != "" {
				// If value is provided, validate it's a proper UID format
				if !rnd.IsUID(it.Value, entity.LabelUID) {
					return fmt.Errorf("invalid label uid format: %s", it.Value)
				}

				labelEntity, err = query.LabelByUID(it.Value)
				if err != nil {
					return fmt.Errorf("label not found: %s", it.Value)
				}
			}
			if labelEntity == nil && it.Title != "" {
				// Create or find by title
				labelEntity = entity.FirstOrCreateLabel(entity.NewLabel(it.Title, 0))
			}

			if labelEntity == nil {
				return fmt.Errorf("could not resolve label to add: value=%s title=%s", it.Value, clean.Log(it.Title))
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
				return fmt.Errorf("label uid required for remove action")
			}

			// Validate UID format
			if !rnd.IsUID(it.Value, entity.LabelUID) {
				return fmt.Errorf("invalid label uid format: %s", it.Value)
			}

			labelEntity, err := query.LabelByUID(it.Value)
			if err != nil || labelEntity == nil || !labelEntity.HasID() {
				return fmt.Errorf("label not found for removal: %s", it.Value)
			}

			if pl, err := query.PhotoLabel(photo.ID, labelEntity.ID); err != nil {
				log.Debugf("batch: photo-label not found for removal: photo=%s label_id=%d", photo.PhotoUID, labelEntity.ID)
				continue
			} else if pl != nil {
				if (pl.LabelSrc == entity.SrcManual || pl.LabelSrc == entity.SrcBatch) && pl.Uncertainty < 100 {
					if err := entity.Db().Delete(&pl).Error; err != nil {
						log.Errorf("batch: delete label failed: %s", err)
					} else {
						log.Debugf("batch: deleted label: photo=%s label_id=%d", photo.PhotoUID, labelEntity.ID)
						changed = true
					}
				} else if pl.LabelSrc != entity.SrcManual && pl.LabelSrc != entity.SrcBatch {
					pl.Uncertainty = 100
					pl.LabelSrc = entity.SrcBatch
					if err := entity.Db().Save(pl).Error; err != nil {
						log.Errorf("batch: block label failed: %s", err)
					} else {
						log.Debugf("batch: blocked label: photo=%s label_id=%d", photo.PhotoUID, labelEntity.ID)
						changed = true
					}
				} else {
					if err := entity.Db().Save(pl).Error; err != nil {
						log.Errorf("batch: save label failed: %s", err)
					}
				}
				_ = photo.RemoveKeyword(labelEntity.LabelName)
				// Persist updated keywords immediately so the change survives reloads
				if err := photo.SaveDetails(); err != nil {
					log.Debugf("batch: failed to save details after keyword removal: %s", err)
				}
			}
		case ActionNone, ActionUpdate:
			// Valid actions that do nothing for labels
			continue
		default:
			return fmt.Errorf("invalid action: %s", it.Action)
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
