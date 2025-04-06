package api

import (
	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/auth/acl"
)

// PostVisionCaption returns a suitable caption for an image.
//
//	@Summary		returns a suitable caption for an image
//	@Id				PostVisionCaption
//	@Tags			Vision
//	@Produce		json
//	@Failure		401,403,404,429,501		{object}	i18n.Response
//	@Router			/api/v1/vision/caption [post]
func PostVisionCaption(router *gin.RouterGroup) {
	router.POST("/vision/caption", func(c *gin.Context) {
		s := Auth(c, acl.ResourceVision, acl.AccessAll)

		// Abort if permission is not granted.
		if s.Abort(c) {
			return
		}

		AbortNotImplemented(c)
	})
}
