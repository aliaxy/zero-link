# Requirements

## Product Scope

zero-link is an administrator-operated short-link system. The first release focuses on reliable link creation, fast redirects, basic analytics, and a simple management UI.

## Functional Requirements

### Administrator Authentication

- Administrators can log in with username and password.
- The system returns an authenticated session token after successful login.
- Management APIs require authentication.
- Public redirect APIs do not require authentication.

### Short-Link Creation

- Administrators can create a short link with an original URL.
- The original URL must use `http` or `https`.
- The system can generate a short code automatically.
- Administrators can provide a custom short code.
- Short codes must be unique.
- Short codes must not use reserved words such as `admin`, `api`, `healthz`, `readyz`, `metrics`, or `static`.
- Administrators can set title, description, enabled state, and expiration time.

### Short-Link Management

- Administrators can view a paginated list of short links.
- The list supports filtering by status and searching by code or title.
- Administrators can view one link's details.
- Administrators can update title, description, original URL, status, and expiration time.
- Short code mutation is not supported after creation.
- Deletion is soft deletion or archival, not physical removal.

### Redirect Resolution

- Visitors can access `GET /{code}`.
- Active and unexpired links return a redirect to the original URL.
- Missing links return `404 Not Found`.
- Disabled links return `403 Forbidden`.
- Expired links return `410 Gone`.
- Deleted links behave like missing links to public visitors.

### Visit Analytics

- The system records visit events after redirect resolution.
- Visit recording must not block successful redirects.
- The system aggregates PV and UV by link and day.
- The management UI can show total PV, total UV, today PV/UV, and recent trends.
- Analytics can be eventually consistent.

### Management UI

- The UI includes login, link list, link creation/editing, link details, and statistics.
- The first UI should be functional and compact, not a complex operations platform.
- The UI uses the management HTTP API and does not bypass service rules.

## Non-Functional Requirements

- Redirect latency should be low and usually served from Redis.
- MySQL is the source of truth for links and statistics.
- Redis is a recoverable cache and rate-limit store.
- Configuration must support local development and Compose deployment.
- Logs must make request failures and dependency failures diagnosable.
- Health checks must distinguish process liveness from dependency readiness.
- Integration tests must be able to run against local MySQL and Redis.

## Boundary Requirements

- Anonymous users cannot create links.
- Internal errors are not exposed directly to visitors or administrators.
- Passwords are never stored in plaintext.
- The system does not trust request headers such as `X-Forwarded-For` unless proxy trust is explicitly configured.
- Statistics are not required to update immediately after each redirect.
