package jwt

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/pkg/fs"
)

var (
	errKeyNotFound = errors.New("jwt: key not found")
)

type cacheEntry struct {
	URL       string      `json:"url"`
	ETag      string      `json:"etag,omitempty"`
	Keys      []PublicJWK `json:"keys"`
	FetchedAt int64       `json:"fetchedAt"`
}

// Verifier validates Portal-issued JWTs on Nodes using JWKS with caching.
type Verifier struct {
	conf *config.Config

	mu        sync.Mutex
	cache     cacheEntry
	cachePath string

	httpClient *http.Client
	now        func() time.Time
}

// ExpectedClaims describes the constraints that must hold for a token.
type ExpectedClaims struct {
	Issuer   string
	Audience string
	Scope    []string
	JWKSURL  string
}

// NewVerifier instantiates a verifier with sane defaults.
func NewVerifier(conf *config.Config) *Verifier {
	v := &Verifier{
		conf:       conf,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		now:        time.Now,
	}
	if conf != nil {
		v.cachePath = filepath.Join(conf.ConfigPath(), "jwks-cache.json")
	}
	_ = v.loadCache()
	return v
}

// Prime ensures JWKS material is cached locally.
func (v *Verifier) Prime(ctx context.Context, jwksURL string) error {
	_, err := v.keysForURL(ctx, jwksURL, true)
	return err
}

// VerifyToken validates a JWT against the expected claims and returns decoded claims.
func (v *Verifier) VerifyToken(ctx context.Context, tokenString string, expected ExpectedClaims) (*Claims, error) {
	if v == nil {
		return nil, errors.New("jwt: verifier not initialized")
	}
	if strings.TrimSpace(tokenString) == "" {
		return nil, errors.New("jwt: token is empty")
	}
	if strings.TrimSpace(expected.Issuer) == "" {
		return nil, errors.New("jwt: expected issuer required")
	}
	if strings.TrimSpace(expected.Audience) == "" {
		return nil, errors.New("jwt: expected audience required")
	}
	if len(expected.Scope) == 0 {
		return nil, errors.New("jwt: expected scope required")
	}

	url := strings.TrimSpace(expected.JWKSURL)
	if url == "" && v.conf != nil {
		url = strings.TrimSpace(v.conf.JWKSUrl())
	}
	if url == "" {
		return nil, errors.New("jwt: jwks url not configured")
	}

	leeway := 60 * time.Second
	if v.conf != nil && v.conf.JWTLeeway() > 0 {
		leeway = time.Duration(v.conf.JWTLeeway()) * time.Second
	}

	parser := gojwt.NewParser(
		gojwt.WithLeeway(leeway),
		gojwt.WithValidMethods([]string{gojwt.SigningMethodEdDSA.Alg()}),
		gojwt.WithIssuer(expected.Issuer),
		gojwt.WithAudience(expected.Audience),
	)

	claims := &Claims{}
	keyFunc := func(token *gojwt.Token) (interface{}, error) {
		kid, _ := token.Header["kid"].(string)
		if kid == "" {
			return nil, errors.New("jwt: missing kid header")
		}
		pk, err := v.publicKeyForKid(ctx, url, kid, false)
		if errors.Is(err, errKeyNotFound) {
			pk, err = v.publicKeyForKid(ctx, url, kid, true)
		}
		if err != nil {
			return nil, err
		}
		return pk, nil
	}

	if _, err := parser.ParseWithClaims(tokenString, claims, keyFunc); err != nil {
		return nil, err
	}

	if claims.IssuedAt == nil || claims.ExpiresAt == nil {
		return nil, errors.New("jwt: missing temporal claims")
	}
	if ttl := claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time); ttl > MaxTokenTTL {
		return nil, errors.New("jwt: token ttl exceeds maximum")
	}

	scopeSet := map[string]struct{}{}
	for _, s := range strings.Fields(claims.Scope) {
		scopeSet[s] = struct{}{}
	}
	for _, req := range expected.Scope {
		if _, ok := scopeSet[req]; !ok {
			return nil, fmt.Errorf("jwt: missing scope %s", req)
		}
	}

	return claims, nil
}

func (v *Verifier) publicKeyForKid(ctx context.Context, url, kid string, force bool) (ed25519.PublicKey, error) {
	keys, err := v.keysForURL(ctx, url, force)
	if err != nil {
		return nil, err
	}
	for _, k := range keys {
		if k.Kid != kid {
			continue
		}
		raw, err := base64.RawURLEncoding.DecodeString(k.X)
		if err != nil {
			return nil, err
		}
		if len(raw) != ed25519.PublicKeySize {
			return nil, fmt.Errorf("jwt: invalid public key length %d", len(raw))
		}
		pk := make(ed25519.PublicKey, ed25519.PublicKeySize)
		copy(pk, raw)
		return pk, nil
	}
	return nil, errKeyNotFound
}

func (v *Verifier) keysForURL(ctx context.Context, url string, force bool) ([]PublicJWK, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	ttl := 300 * time.Second
	if v.conf != nil && v.conf.JWKSCacheTTL() > 0 {
		ttl = time.Duration(v.conf.JWKSCacheTTL()) * time.Second
	}

	if !force && v.cache.URL == url && len(v.cache.Keys) > 0 {
		age := v.now().Unix() - v.cache.FetchedAt
		if age >= 0 && time.Duration(age)*time.Second <= ttl {
			return append([]PublicJWK(nil), v.cache.Keys...), nil
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if v.cache.URL == url && v.cache.ETag != "" {
		req.Header.Set("If-None-Match", v.cache.ETag)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		if v.cache.URL == url && len(v.cache.Keys) > 0 {
			return append([]PublicJWK(nil), v.cache.Keys...), nil
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		v.cache.FetchedAt = v.now().Unix()
		_ = v.saveCacheLocked()
		return append([]PublicJWK(nil), v.cache.Keys...), nil
	}

	if resp.StatusCode != http.StatusOK {
		if v.cache.URL == url && len(v.cache.Keys) > 0 {
			return append([]PublicJWK(nil), v.cache.Keys...), nil
		}
		return nil, fmt.Errorf("jwt: jwks fetch failed: %s", resp.Status)
	}

	var body JWKS
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if len(body.Keys) == 0 {
		return nil, errors.New("jwt: jwks contains no keys")
	}

	v.cache = cacheEntry{
		URL:       url,
		ETag:      resp.Header.Get("ETag"),
		Keys:      append([]PublicJWK(nil), body.Keys...),
		FetchedAt: v.now().Unix(),
	}
	_ = v.saveCacheLocked()

	return append([]PublicJWK(nil), body.Keys...), nil
}

func (v *Verifier) loadCache() error {
	if v.cachePath == "" || !fs.FileExists(v.cachePath) {
		return nil
	}

	b, err := os.ReadFile(v.cachePath)
	if err != nil || len(b) == 0 {
		return err
	}

	var entry cacheEntry
	if err := json.Unmarshal(b, &entry); err != nil {
		return err
	}

	v.cache = entry
	return nil
}

func (v *Verifier) saveCacheLocked() error {
	if v.cachePath == "" {
		return nil
	}
	if err := fs.MkdirAll(filepath.Dir(v.cachePath)); err != nil {
		return err
	}
	data, err := json.Marshal(v.cache)
	if err != nil {
		return err
	}
	return os.WriteFile(v.cachePath, data, fs.ModeSecretFile)
}
