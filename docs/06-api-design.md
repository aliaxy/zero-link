# API Design

## HTTP Conventions

- Management APIs return JSON.
- Redirect APIs return HTTP redirects or simple browser-facing error pages.
- Internal technical errors are logged server-side and exposed as stable public error codes.
- All management APIs except login require administrator authentication.

## Management HTTP API

### POST /admin/login

Request:

```json
{
  "username": "admin",
  "password": "secret"
}
```

Response:

```json
{
  "token": "jwt-token",
  "expires_at": "2026-05-24T12:00:00Z"
}
```

### GET /admin/profile

Returns the authenticated administrator profile.

### POST /admin/links

Request:

```json
{
  "origin_url": "https://example.com/page",
  "code": "campaign1",
  "title": "Campaign 1",
  "description": "Optional note",
  "expire_at": "2026-12-31T23:59:59Z"
}
```

`code` is optional. If omitted, the system generates one.

### GET /admin/links

Query parameters:

- `page`
- `page_size`
- `status`
- `keyword`

Returns paginated short-link summaries.

### GET /admin/links/{id}

Returns one short link and basic aggregate counters.

### PATCH /admin/links/{id}

Updates mutable fields:

- `origin_url`
- `title`
- `description`
- `status`
- `expire_at`

The short code is immutable.

### DELETE /admin/links/{id}

Soft deletes or archives the link.

### GET /admin/links/{id}/stats

Query parameters:

- `from`
- `to`

Returns totals and daily trend data.

## Redirect API

### GET /{code}

Results:

- `302 Found` with `Location` for active links.
- `404 Not Found` for missing or deleted links.
- `403 Forbidden` for disabled links.
- `410 Gone` for expired links.

## Health Checks

### GET /healthz

Returns process liveness. This endpoint should not fail only because MySQL or Redis is temporarily unavailable.

### GET /readyz

Returns readiness. It checks dependencies needed to serve traffic, including RPC, MySQL, and Redis.

## RPC Methods

### CreateShortLink

Creates a link, validates code and URL, and returns the created link.

### GetShortLink

Fetches link details by ID.

### UpdateShortLink

Updates mutable link fields and invalidates cache.

### ResolveShortLink

Resolves a code into redirect data and domain status.

### RecordVisit

Records a visit event. Failure must be observable but should not block redirects.

### GetLinkStats

Returns aggregate statistics for a link and time range.

## Response Envelope

Management APIs should use a stable envelope:

```json
{
  "code": "OK",
  "message": "ok",
  "data": {}
}
```

Error example:

```json
{
  "code": "NOT_FOUND",
  "message": "short link not found"
}
```
