# Security Design

## Trust Boundaries

- Public redirect requests are untrusted.
- Management UI and management API requests are untrusted until authenticated.
- Internal RPC calls are trusted only within the configured deployment network.
- MySQL and Redis contain sensitive operational data and must not be exposed publicly.
- Reverse proxy headers are untrusted unless the proxy is explicitly configured as trusted.

## Authentication

- Administrators authenticate with username and password.
- Passwords are stored using bcrypt or Argon2id.
- Successful login returns a JWT or equivalent signed token.
- Tokens must expire.
- Login failures should return generic messages.

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
