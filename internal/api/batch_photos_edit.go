package api

import (
	"fmt"
	"net/http"

	"github.com/dustin/go-humanize/english"
	"github.com/gin-gonic/gin"
	"github.com/ulule/deepcopier"

	"github.com/photoprism/photoprism/internal/auth/acl"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/entity/query"
	"github.com/photoprism/photoprism/internal/entity/search"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/form/batch"
	"github.com/photoprism/photoprism/internal/photoprism/get"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/i18n"
	"github.com/photoprism/photoprism/pkg/rnd"
)

// BatchPhotosEdit returns and updates the metadata of multiple photos.
//
//	@Summary	returns and updates the metadata of multiple photos
//	@Id			BatchPhotosEdit
//	@Tags		Photos
//	@Accept		json
//	@Produce	json
//	@Success	200						{object}	batch.PhotosResponse
//	@Failure	400,401,403,404,429,500	{object}	i18n.Response
//	@Param		Request					body		batch.PhotosRequest	true	"photos selection and values"
//	@Router		/api/v1/batch/photos/edit [post]
func BatchPhotosEdit(router *gin.RouterGroup) {
	router.Match(MethodsPutPost, "/batch/photos/edit", func(c *gin.Context) {
		s := Auth(c, acl.ResourcePhotos, acl.ActionUpdate)

		if s.Abort(c) {
			return
		}

		conf := get.Config()

		if !conf.Develop() && !conf.Experimental() {
			AbortNotImplemented(c)
			return
		}

		var frm batch.PhotosRequest

		// Assign and validate request form values.
		if err := c.BindJSON(&frm); err != nil {
			AbortBadRequest(c, err)
			return
		}

		if len(frm.Photos) == 0 {
			Abort(c, http.StatusBadRequest, i18n.ErrNoItemsSelected)
			return
		}

		// Fetch selected photos from database.
		photos, count, err := search.BatchPhotos(frm.Photos, s)

		log.Debugf("batch: %s selected for editing", english.Plural(count, "photo", "photos"))

		// Abort if no photos were found.
		if err != nil {
			log.Errorf("batch: %s", clean.Error(err))
			AbortUnexpectedError(c)
			return
		}

		// Update photo metadata based on submitted form values.
		if frm.Values != nil {
			log.Debugf("batch: updating photo metadata for %d photos", len(photos))
			updatedCount := 0

			for i, photo := range photos {
				photoID := photo.PhotoUID

				// Get the full photo entity with preloaded data
				fullPhoto, err := query.PhotoPreloadByUID(photoID)
				if err != nil {
					log.Errorf("batch: failed to load photo %s: %s", photoID, err)
					continue
				}

				// Convert batch form to regular photo form
				photoForm, err := convertBatchToPhotoForm(&fullPhoto, frm.Values)
				if err != nil {
					log.Errorf("batch: failed to convert form for photo %s: %s", photoID, err)
					continue
				}

				// Use the same save mechanism as normal edit
				if err := entity.SavePhotoForm(&fullPhoto, *photoForm); err != nil {
					log.Errorf("batch: failed to save photo %s: %s", photoID, err)
					continue
				}

				// Apply Albums updates if requested
				if frm.Values.Albums.Action == batch.ActionUpdate {
					if err := applyBatchAlbums(photoID, frm.Values.Albums); err != nil {
						log.Errorf("batch: failed to update albums for photo %s: %s", photoID, err)
					}
				}

				// Apply Labels updates if requested
				if frm.Values.Labels.Action == batch.ActionUpdate {
					if err := applyBatchLabels(&fullPhoto, frm.Values.Labels); err != nil {
						log.Errorf("batch: failed to update labels for photo %s: %s", photoID, err)
					}
				}

				// Convert the updated entity.Photo back to search.Photo and update the results array
				updatedSearchPhoto, convertErr := convertEntityToSearchPhoto(&fullPhoto)
				if convertErr != nil {
					log.Errorf("batch: failed to convert photo %s to search result: %s", photoID, convertErr)
				} else {
					photos[i] = *updatedSearchPhoto
				}
				updatedCount++

				// Save sidecar YAML if enabled
				SaveSidecarYaml(&fullPhoto)

				log.Debugf("batch: successfully updated photo %s", photoID)
			}

			log.Infof("batch: successfully updated %d out of %d photos", updatedCount, len(photos))

			// Publish photo update events
			for _, photo := range photos {
				PublishPhotoEvent(StatusUpdated, photo.PhotoUID, c)
			}

			// Update client config and flush cache
			UpdateClientConfig()
			FlushCoverCache()
		}

		// Create batch edit form values form from photo metadata.
		batchFrm := batch.NewPhotosForm(photos)

		// Return models and form values.
		data := batch.PhotosResponse{
			Models: photos,
			Values: batchFrm,
		}

		c.JSON(http.StatusOK, data)
	})
}

// convertBatchToPhotoForm converts batch form values to a regular photo form,
// applying only the fields that have action=update. This allows us to use
// the same SavePhotoForm logic as the normal edit dialog.
func convertBatchToPhotoForm(photo *entity.Photo, batchValues *batch.PhotosForm) (*form.Photo, error) {
	if photo == nil || batchValues == nil {
		return nil, fmt.Errorf("photo or batch values is nil")
	}

	// Start with a form created from the current photo
	photoForm, err := form.NewPhoto(photo)
	if err != nil {
		return nil, fmt.Errorf("failed to create form from photo: %w", err)
	}

	switch batchValues.PhotoTitle.Action {
	case batch.ActionUpdate:
		photoForm.PhotoTitle = batchValues.PhotoTitle.Value
		photoForm.TitleSrc = entity.SrcBatch
	case batch.ActionRemove:
		photoForm.PhotoTitle = ""
		photoForm.TitleSrc = entity.SrcBatch
	}

	switch batchValues.PhotoCaption.Action {
	case batch.ActionUpdate:
		photoForm.PhotoCaption = batchValues.PhotoCaption.Value
		photoForm.CaptionSrc = entity.SrcBatch
	case batch.ActionRemove:
		photoForm.PhotoCaption = ""
		photoForm.CaptionSrc = entity.SrcBatch
	}

	if batchValues.PhotoType.Action == batch.ActionUpdate {
		photoForm.PhotoType = batchValues.PhotoType.Value
		photoForm.TypeSrc = entity.SrcBatch
	}

	// Date/time fields
	timeChanged := false
	if batchValues.PhotoDay.Action == batch.ActionUpdate {
		photoForm.PhotoDay = batchValues.PhotoDay.Value
		timeChanged = true
	}

	if batchValues.PhotoMonth.Action == batch.ActionUpdate {
		photoForm.PhotoMonth = batchValues.PhotoMonth.Value
		timeChanged = true
	}

	if batchValues.PhotoYear.Action == batch.ActionUpdate {
		photoForm.PhotoYear = batchValues.PhotoYear.Value
		timeChanged = true
	}

	if batchValues.TimeZone.Action == batch.ActionUpdate {
		photoForm.TimeZone = batchValues.TimeZone.Value
		timeChanged = true
	}

	if timeChanged {
		photoForm.TakenSrc = entity.SrcBatch
	}

	// Location fields
	locationChanged := false
	if batchValues.PhotoLat.Action == batch.ActionUpdate {
		photoForm.PhotoLat = batchValues.PhotoLat.Value
		locationChanged = true
	}

	if batchValues.PhotoLng.Action == batch.ActionUpdate {
		photoForm.PhotoLng = batchValues.PhotoLng.Value
		locationChanged = true
	}

	if batchValues.PhotoCountry.Action == batch.ActionUpdate {
		photoForm.PhotoCountry = batchValues.PhotoCountry.Value
		locationChanged = true
	}

	if batchValues.PhotoAltitude.Action == batch.ActionUpdate {
		photoForm.PhotoAltitude = batchValues.PhotoAltitude.Value
		locationChanged = true
	}

	if locationChanged {
		photoForm.PlaceSrc = entity.SrcBatch
	}

	// Boolean flags
	if batchValues.PhotoFavorite.Action == batch.ActionUpdate {
		photoForm.PhotoFavorite = batchValues.PhotoFavorite.Value
	}

	if batchValues.PhotoPrivate.Action == batch.ActionUpdate {
		photoForm.PhotoPrivate = batchValues.PhotoPrivate.Value
	}

	if batchValues.PhotoScan.Action == batch.ActionUpdate {
		photoForm.PhotoScan = batchValues.PhotoScan.Value
	}

	if batchValues.PhotoPanorama.Action == batch.ActionUpdate {
		photoForm.PhotoPanorama = batchValues.PhotoPanorama.Value
	}

	// Details fields - preserve existing values, only update changed ones
	currentDetails := photo.GetDetails()
	if currentDetails != nil {
		// Start with current values to preserve unchanged fields
		photoForm.Details.Subject = currentDetails.Subject
		photoForm.Details.SubjectSrc = currentDetails.SubjectSrc
		photoForm.Details.Artist = currentDetails.Artist
		photoForm.Details.ArtistSrc = currentDetails.ArtistSrc
		photoForm.Details.Copyright = currentDetails.Copyright
		photoForm.Details.CopyrightSrc = currentDetails.CopyrightSrc
		photoForm.Details.License = currentDetails.License
		photoForm.Details.LicenseSrc = currentDetails.LicenseSrc
		photoForm.Details.Keywords = currentDetails.Keywords
		photoForm.Details.KeywordsSrc = currentDetails.KeywordsSrc
		photoForm.Details.Notes = currentDetails.Notes
		photoForm.Details.NotesSrc = currentDetails.NotesSrc
	}

	// Now apply only the fields that have action=update
	switch batchValues.DetailsSubject.Action {
	case batch.ActionUpdate:
		photoForm.Details.Subject = batchValues.DetailsSubject.Value
		photoForm.Details.SubjectSrc = entity.SrcBatch
	case batch.ActionRemove:
		photoForm.Details.Subject = ""
		photoForm.Details.SubjectSrc = entity.SrcBatch
	}

	switch batchValues.DetailsArtist.Action {
	case batch.ActionUpdate:
		photoForm.Details.Artist = batchValues.DetailsArtist.Value
		photoForm.Details.ArtistSrc = entity.SrcBatch
	case batch.ActionRemove:
		photoForm.Details.Artist = ""
		photoForm.Details.ArtistSrc = entity.SrcBatch
	}

	switch batchValues.DetailsCopyright.Action {
	case batch.ActionUpdate:
		photoForm.Details.Copyright = batchValues.DetailsCopyright.Value
		photoForm.Details.CopyrightSrc = entity.SrcBatch
	case batch.ActionRemove:
		photoForm.Details.Copyright = ""
		photoForm.Details.CopyrightSrc = entity.SrcBatch
	}

	switch batchValues.DetailsLicense.Action {
	case batch.ActionUpdate:
		photoForm.Details.License = batchValues.DetailsLicense.Value
		photoForm.Details.LicenseSrc = entity.SrcBatch
	case batch.ActionRemove:
		photoForm.Details.License = ""
		photoForm.Details.LicenseSrc = entity.SrcBatch
	}

	// Set the PhotoID for details
	photoForm.Details.PhotoID = photo.ID

	return &photoForm, nil
}

// convertEntityToSearchPhoto converts an entity.Photo to search.Photo for API responses.
func convertEntityToSearchPhoto(photo *entity.Photo) (*search.Photo, error) {
	searchPhoto := &search.Photo{}

	// Copy common fields automatically
	deepcopier.Copy(searchPhoto).From(photo)

	// Set required fields manually
	searchPhoto.CompositeID = fmt.Sprintf("%d", photo.ID)

	// Copy details if they exist
	if details := photo.GetDetails(); details != nil {
		searchPhoto.DetailsSubject = details.Subject
		searchPhoto.DetailsArtist = details.Artist
		searchPhoto.DetailsCopyright = details.Copyright
		searchPhoto.DetailsLicense = details.License
	}

	return searchPhoto, nil
}

// applyBatchAlbums adds/removes the given photo to/from albums according to items action.
func applyBatchAlbums(photoUID string, albums batch.Items) error {
	var addTargets []string

	for _, it := range albums.Items {
		switch it.Action {
		case batch.ActionAdd:
			// Add by UID if provided, otherwise use title to create/find
			if it.Value != "" {
				addTargets = append(addTargets, it.Value)
			} else if it.Title != "" {
				addTargets = append(addTargets, it.Title)
			}
		case batch.ActionRemove:
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

// applyBatchLabels adds/removes labels on the given photo according to items action.
func applyBatchLabels(photo *entity.Photo, labels batch.Items) error {
	if photo == nil || !photo.HasID() {
		return fmt.Errorf("invalid photo")
	}

	// Track if we changed anything to call SaveLabels once
	changed := false

	for _, it := range labels.Items {
		switch it.Action {
		case batch.ActionAdd:
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

		case batch.ActionRemove:
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
