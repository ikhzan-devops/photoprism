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
	"github.com/photoprism/photoprism/internal/photoprism/batch"
	"github.com/photoprism/photoprism/internal/photoprism/get"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/i18n"
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
				photoForm, err := batch.ConvertToPhotoForm(&fullPhoto, frm.Values)
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
					if err := batch.ApplyAlbums(photoID, frm.Values.Albums); err != nil {
						log.Errorf("batch: failed to update albums for photo %s: %s", photoID, err)
					}
				}

				// Apply Labels updates if requested
				if frm.Values.Labels.Action == batch.ActionUpdate {
					if err := batch.ApplyLabels(&fullPhoto, frm.Values.Labels); err != nil {
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
