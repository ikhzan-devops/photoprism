package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/auth/acl"
	clusterjwt "github.com/photoprism/photoprism/internal/auth/jwt"
	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/photoprism/get"
	"github.com/photoprism/photoprism/pkg/authn"
	"github.com/photoprism/photoprism/pkg/clean"
)

// authAnyJWT attempts to authenticate a Portal-issued JWT when a cluster
// node receives a request without an existing session. It verifies the token
// against the node's cached JWKS, ensures the issuer/audience/scope match the
// expected portal values, and, if valid, returns a client session mirroring the
// JWT claims. It returns nil on any validation failure so the caller can fall
// back to existing auth flows. Currently cluster and vision resources are
// eligible for JWT-based authorization; vision access requires the `vision`
// scope whereas cluster access requires the `cluster` scope.
func authAnyJWT(c *gin.Context, clientIP, authToken string, resource acl.Resource, perms acl.Permissions) *entity.Session {
	if c == nil || authToken == "" {
		return nil
	}

	_ = perms

	if resource != acl.ResourceCluster && resource != acl.ResourceVision {
		return nil
	}

	// Basic sanity check for JWT structure.
	if strings.Count(authToken, ".") != 2 {
		return nil
	}

	conf := get.Config()

	if conf == nil || conf.IsPortal() {
		return nil
	}

	if conf.JWKSUrl() == "" {
		return nil
	}

	requiredScopes := []string{"cluster"}
	if resource == acl.ResourceVision {
		requiredScopes = []string{"vision"}
	}

	expected := clusterjwt.ExpectedClaims{
		Audience: fmt.Sprintf("node:%s", conf.NodeUUID()),
		Scope:    requiredScopes,
		JWKSURL:  conf.JWKSUrl(),
	}

	issuers := jwtIssuerCandidates(conf)

	if len(issuers) == 0 {
		return nil
	}

	var (
		claims *clusterjwt.Claims
		err    error
	)

	ctx := c.Request.Context()

	for _, issuer := range issuers {
		expected.Issuer = issuer
		claims, err = get.VerifyJWT(ctx, authToken, expected)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil
	} else if claims == nil {
		return nil
	}

	sess := &entity.Session{
		Status:       http.StatusOK,
		ClientUID:    claims.Subject,
		AuthScope:    clean.Scope(claims.Scope),
		AuthIssuer:   claims.Issuer,
		AuthID:       claims.ID,
		GrantType:    authn.GrantJwtBearer.String(),
		AuthProvider: authn.ProviderClient.String(),
	}

	sess.SetMethod(authn.MethodJWT)
	sess.SetClientName(claims.Subject)
	sess.SetClientIP(clientIP)

	return sess
}

// jwtIssuerCandidates returns the possible issuer values the node should accept
// for Portal JWTs. It prefers the explicit portal cluster identifier and then
// falls back to configured URLs so legacy installations migrate seamlessly.
func jwtIssuerCandidates(conf *config.Config) []string {
	var out []string
	if uuid := conf.ClusterUUID(); uuid != "" {
		out = append(out, fmt.Sprintf("portal:%s", uuid))
	}
	if portal := strings.TrimSpace(conf.PortalUrl()); portal != "" {
		out = append(out, strings.TrimRight(portal, "/"))
	}
	if site := strings.TrimSpace(conf.SiteUrl()); site != "" {
		out = append(out, strings.TrimRight(site, "/"))
	}
	return out
}
