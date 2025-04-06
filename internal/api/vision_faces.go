package api

import (
	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/auth/acl"
)

// PostVisionFaces returns the positions and embeddings of detected faces.
//
//	@Summary		returns the positions and embeddings of detected faces
//	@Id				PostVisionFaces
//	@Tags			Vision
//	@Produce		json
//	@Failure		401,403,429,501		{object}	i18n.Response
//	@Router			/api/v1/vision/faces [post]
func PostVisionFaces(router *gin.RouterGroup) {
	router.POST("/vision/faces", func(c *gin.Context) {
		s := Auth(c, acl.ResourceVision, acl.AccessAll)

		// Abort if permission is not granted.
		if s.Abort(c) {
			return
		}

		AbortNotImplemented(c)
	})
}
