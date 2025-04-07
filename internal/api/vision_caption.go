package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/ai/vision"
	"github.com/photoprism/photoprism/internal/auth/acl"
)

// PostVisionCaption returns a suitable caption for an image.
//
//	@Summary	returns a suitable caption for an image
//	@Id			PostVisionCaption
//	@Tags		Vision
//	@Produce	json
//	@Success	200					{object}	vision.ApiResponse
//	@Failure	401,403,404,429,501	{object}	i18n.Response
//	@Param		images				body		vision.ApiRequest	true	"list of image file urls"
//	@Router		/api/v1/vision/caption [post]
func PostVisionCaption(router *gin.RouterGroup) {
	router.POST("/vision/caption", func(c *gin.Context) {
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
			Model:  &vision.Model{Name: "Caption", Version: "Test", Resolution: 224},
			Result: &vision.ApiResult{Caption: &vision.CaptionResult{Text: "Lorem Ipsum is simply dummy text of the printing and typesetting industry.", Confidence: 0.42424}},
		}

		c.JSON(http.StatusOK, response)
	})
}
