package server

import (
	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/api"
	"github.com/photoprism/photoprism/internal/config"
)

// registerSharingRoutes adds routes for link sharing.
func registerSharingRoutes(router *gin.Engine, conf *config.Config) {
	// Return if the web user interface is disabled.
	if conf.DisableFrontend() {
		return
	}

	s := router.Group(conf.BaseUri("/s"))
	{
		api.Shares(s)
		api.SharePreview(s)
	}
}
