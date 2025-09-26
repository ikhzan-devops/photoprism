package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/photoprism/photoprism/internal/auth/acl"
	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/photoprism/get"
)

func TestAuthAnyJWT(t *testing.T) {
	t.Run("ClusterScope", func(t *testing.T) {
		fx := newPortalJWTFixture(t, "cluster-jwt-success")
		spec := fx.defaultClaimsSpec()
		token := fx.issue(t, spec)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/cluster/theme", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.RemoteAddr = "192.0.2.10:12345"
		c.Request = req

		session := authAnyJWT(c, "192.0.2.10", token, acl.ResourceCluster, nil)
		require.NotNil(t, session)
		assert.Equal(t, http.StatusOK, session.HttpStatus())
		assert.Equal(t, spec.Subject, session.ClientUID)
		assert.Contains(t, session.AuthScope, "cluster")
		assert.Equal(t, spec.Issuer, session.AuthIssuer)
	})

	t.Run("VisionScope", func(t *testing.T) {
		fx := newPortalJWTFixture(t, "cluster-jwt-vision")
		spec := fx.defaultClaimsSpec()
		spec.Scope = []string{"vision"}
		token := fx.issue(t, spec)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/vision/status", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.RemoteAddr = "198.18.0.5:8080"
		c.Request = req

		session := authAnyJWT(c, "198.18.0.5", token, acl.ResourceVision, nil)
		require.NotNil(t, session)
		assert.Equal(t, http.StatusOK, session.HttpStatus())
		assert.Contains(t, session.AuthScope, "vision")
		assert.Equal(t, spec.Issuer, session.AuthIssuer)
	})
	t.Run("RejectsMalformedOrUnknown", func(t *testing.T) {
		fx := newPortalJWTFixture(t, "cluster-jwt-invalid")
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/cluster/theme", nil)
		req.Header.Set("Authorization", "Bearer invalid-token-without-dots")
		req.RemoteAddr = "192.0.2.10:12345"
		c.Request = req

		assert.Nil(t, authAnyJWT(c, "192.0.2.10", "invalid-token-without-dots", acl.ResourceCluster, nil))

		// Ensure we also bail out when JWKS URL is not configured.
		fx.nodeConf.SetJWKSUrl("")
		get.SetConfig(fx.nodeConf)
		assert.Nil(t, authAnyJWT(c, "192.0.2.10", "", acl.ResourceCluster, nil))
	})
	t.Run("NoIssuerMatch", func(t *testing.T) {
		fx := newPortalJWTFixture(t, "cluster-jwt-no-issuer")
		spec := fx.defaultClaimsSpec()
		token := fx.issue(t, spec)

		// Remove all issuer candidates.
		origPortal := fx.nodeConf.Options().PortalUrl
		origSite := fx.nodeConf.Options().SiteUrl
		origClusterUUID := fx.nodeConf.Options().ClusterUUID
		fx.nodeConf.Options().PortalUrl = ""
		fx.nodeConf.Options().SiteUrl = ""
		fx.nodeConf.Options().ClusterUUID = ""
		get.SetConfig(fx.nodeConf)
		t.Cleanup(func() {
			fx.nodeConf.Options().PortalUrl = origPortal
			fx.nodeConf.Options().SiteUrl = origSite
			fx.nodeConf.Options().ClusterUUID = origClusterUUID
			get.SetConfig(fx.nodeConf)
		})

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/cluster/theme", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.RemoteAddr = "203.0.113.5:2222"
		c.Request = req

		assert.Nil(t, authAnyJWT(c, "203.0.113.5", token, acl.ResourceCluster, nil))
	})
	t.Run("UnsupportedResource", func(t *testing.T) {
		fx := newPortalJWTFixture(t, "cluster-jwt-unsupported")
		token := fx.issue(t, fx.defaultClaimsSpec())

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/cluster/theme", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.RemoteAddr = "198.51.100.7:9999"
		c.Request = req

		assert.Nil(t, authAnyJWT(c, "198.51.100.7", token, acl.ResourcePhotos, nil))
	})
}

func TestJwtIssuerCandidates(t *testing.T) {
	t.Run("IncludesAllSources", func(t *testing.T) {
		conf := config.NewConfig(config.CliTestContext())
		conf.Options().ClusterUUID = "11111111-1111-4111-8111-111111111111"
		conf.Options().PortalUrl = "https://portal.example.test/"
		conf.Options().SiteUrl = "https://site.example.test/base/"

		orig := get.Config()
		get.SetConfig(conf)
		t.Cleanup(func() { get.SetConfig(orig) })

		cands := jwtIssuerCandidates(conf)
		assert.Equal(t, []string{
			"portal:11111111-1111-4111-8111-111111111111",
			"https://portal.example.test",
			"https://site.example.test/base",
		}, cands)
	})
	t.Run("DefaultsToLocalhost", func(t *testing.T) {
		conf := config.NewConfig(config.CliTestContext())
		conf.Options().ClusterUUID = ""
		conf.Options().PortalUrl = ""
		conf.Options().SiteUrl = ""

		assert.Equal(t, []string{"http://localhost:2342"}, jwtIssuerCandidates(conf))
	})
}
