package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/api/hooks"
	"github.com/photoprism/photoprism/internal/auth/acl"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/internal/photoprism/get"
	"github.com/photoprism/photoprism/pkg/authn"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/i18n"
	"github.com/photoprism/photoprism/pkg/media/http/header"
)

// Webhook handles web event hooks and checks their authorization.
//
//	@Summary	handles web event hooks and checks their authorization
//	@Id			Webhook
//	@Tags		Webhook
//	@Accept		json
//	@Success	200
//	@Failure	401,403,429
//	@Param		payload	body	hooks.Payload	true	"webhook event data"
//	@Router		/api/v1/webhook/{channel} [post]
func Webhook(router *gin.RouterGroup) {
	instanceRequestHandler := func(c *gin.Context) {
		s := Auth(c, acl.ResourceWebhooks, acl.ActionPublish)

		if s.Abort(c) {
			return
		}

		var request hooks.Payload

		// Assign and validate request form values.
		if c.Request.Method == http.MethodGet {
			if err := c.BindQuery(&request); err != nil {
				event.AuditErr([]string{ClientIP(c), "session %s", "webhook", "%s"}, s.RefID, err)
				AbortBadRequest(c)
				return
			}
		} else {
			if err := c.BindJSON(&request); err != nil {
				event.AuditErr([]string{ClientIP(c), "session %s", "webhook", "%s"}, s.RefID, err)
				AbortBadRequest(c)
				return
			}
		}

		eventType := clean.TypeLowerUnderscore(request.Type)

		if eventType == "" {
			event.AuditWarn([]string{ClientIP(c), "session %s", "webhook", "missing type"}, s.RefID)
			AbortBadRequest(c)
			return
		}

		if request.Data == nil {
			event.AuditWarn([]string{ClientIP(c), "session %s", "webhook", "missing data"}, s.RefID)
			AbortBadRequest(c)
			return
		}

		resource, resourceEv, found := strings.Cut(eventType, ".")

		if !found || resource == "" || resourceEv == "" {
			event.AuditWarn([]string{ClientIP(c), "session %s", "webhook", "%s", authn.Denied}, s.RefID, eventType)
			AbortBadRequest(c)
			return
		}

		ev := "instance." + eventType

		if s.IsClient() {
			if acl.Rules.Deny(acl.Resource(resource), s.ClientRole(), acl.ActionPublish) {
				event.AuditWarn([]string{ClientIP(c), "client %s", "session %s", "webhook", "%s", authn.Denied}, clean.Log(s.ClientInfo()), s.RefID, ev)
				AbortForbidden(c)
				return
			} else if s.InsufficientScope(acl.Resource(resource), acl.Permissions{acl.ActionPublish}) {
				event.AuditErr([]string{ClientIP(c), "client %s", "session %s", "webhook", "%s", authn.ErrInsufficientScope.Error()}, clean.Log(s.ClientInfo()), s.RefID, ev)
				AbortForbidden(c)
				return
			}
		} else {
			event.AuditWarn([]string{ClientIP(c), "session %s", "webhook", "%s", authn.Denied}, s.RefID, eventType)
			AbortForbidden(c)
			return
		}

		switch ev {
		// case "instance.api.downloads.register":
		//	_ = download.Register(fmt.Sprintf("%v", request.Data["uuid"]), fmt.Sprintf("%v", request.Data["filename"]))
		default:
			event.Publish(ev, request.Data)
		}

	}

	requestHandler := func(c *gin.Context) {
		// Prevent API response caching.
		c.Header(header.CacheControl, header.CacheControlNoStore)

		// For security reasons, this endpoint is not available in public or demo mode.
		if conf := get.Config(); conf.Public() || conf.Demo() {
			Abort(c, http.StatusForbidden, i18n.ErrFeatureDisabled)
			return
		}

		// Only the instance channel is currently implemented.
		channel := clean.Token(c.Param("channel"))

		switch channel {
		case acl.ChannelInstance.String():
			instanceRequestHandler(c)
		default:
			AbortNotImplemented(c)
			return
		}
	}

	router.GET("/webhook/:channel", requestHandler)
	router.POST("/webhook/:channel", requestHandler)
}
