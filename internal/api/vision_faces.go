package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/ai/vision"
	"github.com/photoprism/photoprism/internal/auth/acl"
)

// PostVisionFaces returns the positions and embeddings of detected faces.
//
//	@Summary	returns the positions and embeddings of detected faces
//	@Id			PostVisionFaces
//	@Tags		Vision
//	@Produce	json
//	@Success	200				{object}	vision.ApiResponse
//	@Failure	401,403,429,501	{object}	i18n.Response
//	@Param		images			body		vision.ApiRequest	true	"list of image file urls"
//	@Router		/api/v1/vision/faces [post]
func PostVisionFaces(router *gin.RouterGroup) {
	router.POST("/vision/faces", func(c *gin.Context) {
		s := Auth(c, acl.ResourceVision, acl.AccessAll)

		// Abort if permission is not granted.
		if s.Abort(c) {
			return
		}

		var request vision.ApiRequest

		// Assign and validate request form values.
		if err := c.BindJSON(&request); err != nil {
			AbortBadRequest(c)
			return
		}

		// Generate Vision API service response.
		response := vision.ApiResponse{
			Id:     request.GetId(),
			Model:  &vision.Model{Name: "Faces", Version: "Test", Resolution: 224},
			Result: &vision.ApiResult{Faces: &[]string{}},
		}

		c.JSON(http.StatusOK, response)
	})
}
