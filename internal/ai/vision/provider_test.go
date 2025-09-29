package vision

import "testing"

func TestRegisterProviderAlias(t *testing.T) {
	const alias = "unit-test"
	providerMu.Lock()
	prev, had := providerAliasIndex[alias]
	if had {
		delete(providerAliasIndex, alias)
	}
	providerMu.Unlock()

	t.Cleanup(func() {
		providerMu.Lock()
		if had {
			providerAliasIndex[alias] = prev
		} else {
			delete(providerAliasIndex, alias)
		}
		providerMu.Unlock()
	})

	RegisterProviderAlias("  Unit-Test  ", ProviderInfo{RequestFormat: ApiFormat("custom"), ResponseFormat: "", FileScheme: "data"})

	info, ok := ProviderInfoFor(alias)
	if !ok {
		t.Fatalf("expected provider alias %q to be registered", alias)
	}

	if info.RequestFormat != ApiFormat("custom") {
		t.Errorf("unexpected request format: %s", info.RequestFormat)
	}

	if info.ResponseFormat != ApiFormat("custom") {
		t.Errorf("expected response format default to request, got %s", info.ResponseFormat)
	}

	if info.FileScheme != "data" {
		t.Errorf("unexpected file scheme: %s", info.FileScheme)
	}
}

func TestRegisterProvider(t *testing.T) {
	format := ApiFormat("unit-format")
	provider := Provider{}

	providerMu.Lock()
	prev, had := providerRegistry[format]
	if had {
		delete(providerRegistry, format)
	}
	providerMu.Unlock()

	t.Cleanup(func() {
		providerMu.Lock()
		if had {
			providerRegistry[format] = prev
		} else {
			delete(providerRegistry, format)
		}
		providerMu.Unlock()
	})

	RegisterProvider(format, provider)
	got, ok := ProviderFor(format)
	if !ok {
		t.Fatalf("expected provider for %s", format)
	}

	if got != provider {
		t.Errorf("unexpected provider value: %#v", got)
	}
}
