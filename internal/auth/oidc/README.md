## PhotoPrism — OIDC Integration

**Last Updated:** November 23, 2025

### Overview

`internal/auth/oidc` implements PhotoPrism’s OpenID Connect (OIDC) Relying Party (RP) flow so users can sign in with third‑party identity providers. The package wraps the `zitadel/oidc` client to perform discovery, build the RP, redirect users to the provider, exchange codes for tokens, and retrieve profile claims in a predictable, testable way.

#### Context & Constraints

- Relies on the provider’s `/.well-known/openid-configuration` for discovery and enforces `https` unless explicitly allowed via `insecure`.
- Uses random per-session cookie keys (16‑byte hash + encrypt) and the shared HTTP client defined in `http_client.go`.
- PKCE is enabled automatically when the provider advertises `S256`.
- Scopes default to `authn.OidcRequiredScopes` when none are supplied; scopes are cleaned via `clean.Scopes`.
- Token exchange uses the provider’s userinfo endpoint by default; errors are surfaced via Gin response headers (`oidc_error`) and audit logs.

#### Goals

- Provide a consistent RP client that can be reused by CLI, server routes, and tests.
- Keep redirect and code‑exchange handlers minimal while ensuring audit visibility and secure defaults.
- Allow editions (CE/Pro) to extend claim processing (e.g., groups, roles) without duplicating RP wiring.

#### Non-Goals

- Managing upstream identity provider configuration or enrollment.
- Implementing a full OIDC Provider; PhotoPrism acts only as a Relying Party.
- Handling every custom claim set; extension hooks should live beside claim parsing code.

### Package Layout (Code Map)

- `oidc.go` — package doc + logger.
- `client.go` — RP construction (`NewClient`), PKCE detection, auth redirect, code exchange + userinfo retrieval.
- `http_client.go` — shared HTTP client with TLS toggle and timeouts; helpers for tests in `http_client_test.go`.
- `redirect_url.go` — builds the redirect/callback URL from site config.
- `register.go` — provider registration glue; tests in `register_test.go`.
- `username.go` — derives usernames from claims; tests in `username_test.go`.
- `client_test.go`, `oidc_test.go` — happy-path and error-path coverage for discovery, auth URL, and code exchange.

### Related Packages & Entry Points

- `internal/server/routes.go` registers the OIDC auth and callback endpoints.
- `pkg/authn` defines required scopes and shared auth helpers.
- `internal/auth/acl` and (Pro) `pro/internal/auth/ldap` handle role/group mapping; the planned OIDC group parsing will mirror this logic.
- `internal/config` provides OIDC options/flags (issuer, client ID/secret, scopes, insecure).
- `internal/event` supplies the logger used for audit and error reporting.

### Operational Tips

- Always call `RedirectURL(siteUrl)` to build callbacks that respect reverse proxies and base URIs.
- Reuse `HttpClient(insecure)` so timeouts and TLS settings stay consistent.
- When adding claims processing, keep parsing isolated (e.g., new helper) and ensure failures do not block sign‑in unless required.

### Configuration & Safety

- Enforce `https` for issuers unless `insecure` is explicitly set (intended for dev/test).
- Cookie handler is created per client with fresh random keys to avoid reuse across restarts.
- Audit every provider/redirect/token error with sanitized messages; avoid logging secrets.
- Prefer explicit scopes from configuration; defaults request only the minimal set.

### Security Group Extension for Entra ID

- [x] Parse `groups` claim (string GUIDs) from ID/Access tokens and map to roles using a configurable mapping, mirroring LDAP group handling.
- [x] Detect `_claim_names.groups` overage markers; optionally fetch memberships from Microsoft Graph (`/me/transitiveMemberOf`) behind a feature flag and enforced timeouts.
- [x] Keep `roles`/`wids` (app and directory roles) separate from security groups; document precedence.
- [x] Add unit tests with signed JWT fixtures covering: direct `groups` array, overage marker, and no-groups fallback.
- [x] Expose config knobs (e.g., `AuthOIDCGroupsToRoles`, Graph fetch enable/disable, request timeout) and surface them in CLI/`options.yml` reports.
- [ ] Add integration doc/tests for Entra app registration requirements (`groupMembershipClaims=SecurityGroup|All|ApplicationGroup`) and token size limits (~200 groups).
- [ ] Update Pro parity notes so LDAP and OIDC group mappings share helpers and behavior.

#### Related Resources & Specs

- Microsoft Entra group claims: https://learn.microsoft.com/en-us/entra/identity-platform/access-token-claims-reference#groups-claim
- Group overage handling: https://learn.microsoft.com/en-us/entra/identity-platform/howto-add-app-roles-in-azure-ad-apps#group-overage-and-_claim_names
- Token customization guidance: https://learn.microsoft.com/en-us/entra/architecture/customize-tokens

> **Note:** Entra ID security groups are only supported in PhotoPrism® Pro.

### Testing

- Unit tests: `go test ./internal/auth/oidc -count=1`
- Tests cover discovery failures, PKCE detection, redirect URL construction, username extraction, and code‑exchange error handling.
- For integration testing with a real IdP, set OIDC env vars in `compose.local.yaml`, start the dev server, and exercise `/auth/oidc` + callback.
