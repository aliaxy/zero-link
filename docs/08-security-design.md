# Security Design

## Trust Boundaries

- Public redirect requests are untrusted.
- Management UI and management API requests are untrusted until authenticated.
- Internal RPC calls are trusted only within the configured deployment network.
- MySQL and Redis contain sensitive operational data and must not be exposed publicly.
- Reverse proxy headers are untrusted unless the proxy is explicitly configured as trusted.

## Authentication

- Administrators authenticate with username and password.
- Passwords are stored using bcrypt (`bcrypt.DefaultCost`).
- Successful login returns a dual-token pair:
  - **Access token** — short-lived JWT (≤ 15 min). Stateless; validated locally without a Redis lookup.
    Cannot be revoked before expiry; the short TTL is the accepted trade-off.
  - **Refresh token** — opaque 32-byte `crypto/rand` token (7-day TTL). Only the SHA-256 hash is stored
    in Redis (`zl:rt:{hash}`). Rotation deletes the old hash and issues a new token; submitting an
    already-rotated token is treated as a reuse attempt and returns `UNAUTHENTICATED`.
- `POST /admin/refresh` rotates the refresh token and issues a new access token.
- `PATCH /admin/password` changes the password and calls `RevokeAll`, which bulk-deletes all refresh
  token hashes for the account from Redis. All sessions must re-authenticate; in-flight access tokens
  remain valid until their natural expiry.
- Login failures return generic messages.
- Refresh token raw values are never logged or stored; only the SHA-256 hash persists in Redis.

## Startup Config Validation

Service configuration is validated at startup via `go-playground/validator/v10`. Invalid configuration causes an immediate fatal log and process exit rather than allowing the service to run with a misconfigured state.

Enforced constraints:

- `Auth.Secret`: required, minimum 32 bytes. A short secret reduces effective brute-force resistance for HS256 tokens.
- `Auth.TokenTTLSeconds`: required, greater than zero.
- `Analytics.IPSalt` (link-rpc): required, non-empty. An absent salt would collapse all IP hashes to the same bucket, removing the pseudonymisation guarantee.

## Authorization

- All management endpoints require administrator identity.
- Public redirect endpoints do not require login.
- Future multi-admin support should add explicit role checks before adding new administrator operations.

## Input Safety

### URL Validation

- Only `http` and `https` schemes are allowed.
- Empty hosts are rejected.
- Dangerous schemes such as `javascript`, `file`, and `data` are rejected.
- Optional domain denylist can block known malicious targets.

### Short-Code Validation

- Use an allowlist of characters: letters, digits, hyphen, and underscore.
- Enforce minimum and maximum length.
- Reject reserved words such as `admin`, `api`, `healthz`, `readyz`, `metrics`, and `static`.
- Treat codes case-sensitively only if documented; otherwise normalize consistently.

### Pagination

- Enforce a maximum `page_size`.
- Reject negative page values.

## Abuse Prevention

- Rate limit link creation per administrator.
- `GET /{code}`: per-IP rate limit of 20 req/s using go-zero `PeriodLimit` backed by Redis (`rl:redirect:ip:{ip}`). Excess requests return `429 Too Many Requests`.
- `POST /admin/login`: per-IP rate limit of 10 req/min using go-zero `PeriodLimit` (`rl:login:ip:{ip}`). Excess requests return `429 Too Many Requests`.
- Client IP is extracted from `X-Forwarded-For` (first value) with fallback to `RemoteAddr`; trust only when a reverse proxy is present.
- Consider bot detection later if analytics quality becomes important.
- Avoid negative-cache poisoning by keeping missing-link caches short if introduced.
- Short codes are permanently reserved in `reserved_code` after archival. `CreateShortLink` checks this table before allowing a custom code, preventing an attacker from claiming a previously popular short code after its original link has been archived and hard-deleted.

## Logging Safety

- Do not log plaintext passwords.
- Do not log complete tokens.
- Do not log full IP addresses by default.
- Keep stable log messages low-cardinality; attach IDs and codes as structured fields.
- Avoid logging full original URLs if they may contain sensitive query parameters.

## Web Security

- Management UI should use HttpOnly and SameSite cookies if cookie-based auth is chosen.
- Set basic security headers for browser routes.
- Do not rely on client-side authorization.
- Escape all rendered user-controlled text.

## Dependency Safety

- Prefer standard library and go-zero primitives when adequate.
- Add dependencies only when they remove real complexity.
- Run vulnerability scanning in CI once the Go module exists.

## Open Redirect Consideration

This product intentionally redirects to administrator-provided URLs. The defense is not to prohibit all external redirects, but to authenticate creation, validate URL schemes, optionally deny risky domains, and clearly separate public redirect behavior from management auth flows.
