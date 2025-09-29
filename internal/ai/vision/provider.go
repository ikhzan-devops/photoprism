package vision

import (
	"context"
	"strings"
	"sync"
)

// ModelProvider represents the canonical identifier for a computer vision service provider.
type ModelProvider = string

const (
	// ProviderVision represents the default PhotoPrism vision service endpoints.
	ProviderVision ModelProvider = "vision"
	// ProviderTensorFlow represents on-device TensorFlow models.
	ProviderTensorFlow ModelProvider = "tensorflow"
	// ProviderLocal is used when no explicit provider can be determined.
	ProviderLocal ModelProvider = "local"
)

// RequestBuilder builds an API request for a provider based on the model configuration and input files.
type RequestBuilder interface {
	Build(ctx context.Context, model *Model, files Files) (*ApiRequest, error)
}

// ResponseParser parses a raw provider response into the generic ApiResponse structure.
type ResponseParser interface {
	Parse(ctx context.Context, req *ApiRequest, raw []byte, status int) (*ApiResponse, error)
}

// DefaultsProvider supplies provider-specific prompt and schema defaults when they are not configured explicitly.
type DefaultsProvider interface {
	SystemPrompt(model *Model) string
	UserPrompt(model *Model) string
	SchemaTemplate(model *Model) string
}

// Provider groups the callbacks required to integrate a third-party vision service.
type Provider struct {
	Builder  RequestBuilder
	Parser   ResponseParser
	Defaults DefaultsProvider
}

var (
	providerRegistry   = make(map[ApiFormat]Provider)
	providerAliasIndex = make(map[string]ProviderInfo)
	providerMu         sync.RWMutex
)

// RegisterProvider adds/overrides a provider implementation for a specific API format.
func RegisterProvider(format ApiFormat, provider Provider) {
	providerMu.Lock()
	defer providerMu.Unlock()
	providerRegistry[format] = provider
}

// ProviderInfo describes metadata that can be associated with a provider alias.
type ProviderInfo struct {
	RequestFormat  ApiFormat
	ResponseFormat ApiFormat
	FileScheme     string
}

// RegisterProviderAlias maps a logical provider name (e.g. "ollama") to a request/response format pair.
func RegisterProviderAlias(name string, info ProviderInfo) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" || info.RequestFormat == "" {
		return
	}

	if info.ResponseFormat == "" {
		info.ResponseFormat = info.RequestFormat
	}

	providerMu.Lock()
	providerAliasIndex[name] = info
	providerMu.Unlock()
}

// ProviderInfoFor returns the metadata associated with a logical provider name.
func ProviderInfoFor(name string) (ProviderInfo, bool) {
	name = strings.TrimSpace(strings.ToLower(name))
	providerMu.RLock()
	info, ok := providerAliasIndex[name]
	providerMu.RUnlock()
	return info, ok
}

// ProviderFor returns the registered provider implementation for the given API format, if any.
func ProviderFor(format ApiFormat) (Provider, bool) {
	providerMu.RLock()
	defer providerMu.RUnlock()
	provider, ok := providerRegistry[format]
	return provider, ok
}
