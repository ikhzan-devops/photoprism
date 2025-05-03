package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/auth/acl"
	"github.com/photoprism/photoprism/internal/entity/query"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/form/batch"
	"github.com/photoprism/photoprism/internal/photoprism/get"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/i18n"
)

// BatchPhotos returns the metadata of multiple pictures so that it can be edited.
//
//	@Summary	returns the metadata of multiple pictures so that it can be edited
//	@Id			BatchPhotos
//	@Tags		Photos
//	@Accept		json
//	@Produce	json
//	@Success	200						{object}	batch.PhotoForm
//	@Failure	400,401,403,404,429,500	{object}	i18n.Response
//	@Param		photos					body		form.Selection	true	"Photo Selection"
//	@Router		/api/v1/batch/photos [post]
func BatchPhotos(router *gin.RouterGroup) {
	router.POST("/batch/photos", func(c *gin.Context) {
		s := Auth(c, acl.ResourcePhotos, acl.ActionUpdate)

		if s.Abort(c) {
			return
		}

		conf := get.Config()

		if !conf.Develop() && !conf.Experimental() {
			AbortNotImplemented(c)
			return
		}

		var frm form.Selection

		// Assign and validate request form values.
		if err := c.BindJSON(&frm); err != nil {
			AbortBadRequest(c)
			return
		}

		if len(frm.Photos) == 0 {
			Abort(c, http.StatusBadRequest, i18n.ErrNoItemsSelected)
			return
		}

		// Find selected photos.
		photos, err := query.SelectedPhotos(frm)

		if err != nil {
			log.Errorf("batch: %s", clean.Error(err))
			AbortUnexpectedError(c)
			return
		}

		// Load files and details.
		for _, photo := range photos {
			photo.PreloadFiles()
			photo.GetDetails()
		}

		batchFrm := batch.NewPhotoForm(photos)

		data := gin.H{
			"photos": photos,
			"values": batchFrm,
		}

		c.JSON(http.StatusOK, data)
	})
}
