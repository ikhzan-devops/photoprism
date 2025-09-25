package get

import (
	"sync"

	"github.com/photoprism/photoprism/internal/auth/jwt"
	"github.com/photoprism/photoprism/pkg/clean"
)

var (
	onceJWTManager sync.Once
	onceJWTIssuer  sync.Once
)

func initJWTManager() {
	conf := Config()
	if conf == nil || !conf.IsPortal() {
		return
	}
	manager, err := jwt.NewManager(conf)
	if err != nil {
		log.Warnf("jwt: manager init failed (%s)", clean.Error(err))
		return
	}
	if _, err := manager.EnsureActiveKey(); err != nil {
		log.Warnf("jwt: ensure signing key failed (%s)", clean.Error(err))
	}
	services.JWTManager = manager
}

// JWTManager returns the portal key manager; nil on nodes.
func JWTManager() *jwt.Manager {
	onceJWTManager.Do(initJWTManager)
	return services.JWTManager
}

func initJWTIssuer() {
	manager := JWTManager()
	if manager == nil {
		return
	}
	services.JWTIssuer = jwt.NewIssuer(manager)
}

// JWTIssuer returns the portal JWT issuer helper; nil on nodes.
func JWTIssuer() *jwt.Issuer {
	onceJWTIssuer.Do(initJWTIssuer)
	return services.JWTIssuer
}

// JWTVerifier returns a verifier bound to the current config.
func JWTVerifier() *jwt.Verifier {
	conf := Config()
	if conf == nil {
		return nil
	}
	return jwt.NewVerifier(conf)
}
