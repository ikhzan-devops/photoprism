package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/ai/vision"
	"github.com/photoprism/photoprism/internal/auth/acl"
	"github.com/photoprism/photoprism/pkg/rnd"
)

// PostVisionLabels returns suitable labels for an image.
//
//	@Summary		returns suitable labels for an image
//	@Id				PostVisionLabels
//	@Tags			Vision
//	@Produce		json
//	@Success		200				{object}	vision.LabelsResponse
//	@Failure		401,403,429		{object}	i18n.Response
//	@Router			/api/v1/vision/labels [post]
func PostVisionLabels(router *gin.RouterGroup) {
	router.POST("/vision/labels", func(c *gin.Context) {
		s := Auth(c, acl.ResourceVision, acl.AccessAll)

		// Abort if permission is not granted.
		if s.Abort(c) {
			return
		}

		response := vision.NewLabelsResponse(rnd.UUID(), vision.NasnetModel, nil)

		c.JSON(http.StatusOK, response)
	})
}
