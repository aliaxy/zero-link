# API Design

## HTTP Conventions

- Management APIs return JSON.
- Redirect APIs return HTTP redirects or simple browser-facing error pages.
- Internal technical errors are logged server-side and exposed as stable public error codes.
- Stage 3 management APIs use `Authorization: Bearer <jwt>` authentication.
- All Stage 3 management APIs except login require administrator authentication.
- Stage 3 does not add redirect serving, Redis redirect cache behavior, analytics APIs, or the management UI.

## Response Envelope

Management APIs should use a stable envelope for both success and errors.

Success:

```json
{
  "code": "OK",
  "message": "ok",
  "data": {}
}
```

Error:

```json
{
  "code": "NOT_FOUND",
  "message": "short link not found"
}
```

Common public error codes:

- `OK`: request completed successfully.
- `INVALID_ARGUMENT`: request body, path parameter, or query parameter is invalid.
- `UNAUTHENTICATED`: token is missing, malformed, expired, or invalid.
- `PERMISSION_DENIED`: authenticated administrator is not allowed to perform the action.
- `NOT_FOUND`: requested administrator or short link does not exist, is disabled for auth, or is soft deleted.
- `CONFLICT`: custom short code already exists.
- `INTERNAL`: unexpected server-side failure.

## Authentication And Token Configuration

Stage 3 administrator authentication uses JWT bearer tokens issued by `link-api` after `link-rpc` authenticates administrator credentials.

Configuration fields:

- `Auth.Secret`: local service signing secret for management JWTs.
- `Auth.TokenTTLSeconds`: token lifetime in seconds.

Configuration rules:

- Example configuration documents the fields with non-secret placeholder values.
- Local ignored configuration stores machine-local development values.
- Real secrets must not be committed.
- Token claims include the administrator ID and username.
- Management handlers must reject missing, malformed, expired, or invalid bearer tokens before calling business RPC methods.

## Stage 3 Management HTTP API

### POST /admin/login

Authenticates an administrator and returns a management token.

Authentication:

- Not required.

Request:

```json
{
  "username": "admin",
  "password": "secret"
}
```

Response data:

```json
{
  "token": "jwt-token",
  "expires_at": "2026-12-31T12:00:00Z",
  "admin": {
    "id": 1,
    "username": "admin"
  }
}
```

Errors:

- `INVALID_ARGUMENT` for missing username or password.
- `UNAUTHENTICATED` for invalid credentials or inactive administrators.
- `INTERNAL` for unexpected authentication failures.

### GET /admin/profile

Returns the authenticated administrator profile.

Authentication:

- Required.
- The request must include `Authorization: Bearer <jwt>`.

Response data:

```json
{
  "id": 1,
  "username": "admin",
  "status": 1,
  "created_at": "2026-05-27T12:00:00Z"
}
```

Errors:

- `UNAUTHENTICATED` for missing or invalid bearer token.
- `NOT_FOUND` for a missing or inactive administrator.

### POST /admin/links

Creates a short link for management use.

Authentication:

- Required.

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

Validation:

- `origin_url` is required and must be an absolute `http` or `https` URL.
- `code` is optional.
- When `code` is omitted, the system generates a 6-character base62 code using `crypto/rand`.
- Custom codes must be 3-32 characters.
- Custom codes may contain ASCII letters, digits, `_`, and `-`.
- Reserved words are rejected so management and future system routes cannot be shadowed.
- `title` and `description` are optional.
- `expire_at` is optional; when present, it must be in the future.

Response data:

```json
{
  "id": 1001,
  "code": "campaign1",
  "origin_url": "https://example.com/page",
  "title": "Campaign 1",
  "description": "Optional note",
  "status": 1,
  "expire_at": "2026-12-31T23:59:59Z",
  "created_by": 1,
  "created_at": "2026-05-27T12:00:00Z",
  "updated_at": "2026-05-27T12:00:00Z"
}
```

Errors:

- `INVALID_ARGUMENT` for invalid URL, invalid code, reserved code, invalid status, or invalid expiration.
- `CONFLICT` when a custom code already exists.
- `UNAUTHENTICATED` for missing or invalid bearer token.

### GET /admin/links

Returns paginated short-link summaries.

Authentication:

- Required.

Query parameters:

- `page`: optional positive integer, defaults to `1`.
- `page_size`: optional positive integer, defaults to `20`, maximum `100`.
- `status`: optional status filter.
- `keyword`: optional search term matched against code, title, or origin URL.

Response data:

```json
{
  "items": [
    {
      "id": 1001,
      "code": "campaign1",
      "origin_url": "https://example.com/page",
      "title": "Campaign 1",
      "status": 1,
      "expire_at": "2026-12-31T23:59:59Z",
      "created_at": "2026-05-27T12:00:00Z",
      "updated_at": "2026-05-27T12:00:00Z"
    }
  ],
  "page": 1,
  "page_size": 20,
  "total": 1
}
```

Errors:

- `INVALID_ARGUMENT` for invalid pagination or filter values.
- `UNAUTHENTICATED` for missing or invalid bearer token.

### GET /admin/links/{id}

Returns one short link by ID.

Authentication:

- Required.

Response data:

```json
{
  "id": 1001,
  "code": "campaign1",
  "origin_url": "https://example.com/page",
  "title": "Campaign 1",
  "description": "Optional note",
  "status": 1,
  "expire_at": "2026-12-31T23:59:59Z",
  "created_by": 1,
  "created_at": "2026-05-27T12:00:00Z",
  "updated_at": "2026-05-27T12:00:00Z"
}
```

Errors:

- `INVALID_ARGUMENT` for invalid ID.
- `NOT_FOUND` for missing or soft-deleted links.
- `UNAUTHENTICATED` for missing or invalid bearer token.

### PATCH /admin/links/{id}

Updates mutable short-link fields.

Authentication:

- Required.

Request:

```json
{
  "origin_url": "https://example.com/updated-page",
  "title": "Updated title",
  "description": "Updated note",
  "status": 1,
  "expire_at": "2026-12-31T23:59:59Z"
}
```

Rules:

- `code` is immutable and must not be accepted as an update field.
- Mutable fields are `origin_url`, `title`, `description`, `status`, and `expire_at`.
- Updating a link only changes MySQL state in Stage 3.
- Redirect cache invalidation is documented for Stage 4 and must not be implemented in Stage 3.

Errors:

- `INVALID_ARGUMENT` for invalid ID or invalid mutable field values.
- `NOT_FOUND` for missing or soft-deleted links.
- `UNAUTHENTICATED` for missing or invalid bearer token.

### DELETE /admin/links/{id}

Soft deletes a short link.

Authentication:

- Required.

Response data:

```json
{
  "id": 1001,
  "deleted": true
}
```

Rules:

- Deletion sets `deleted_at`.
- Soft-deleted links are hidden from normal management list and detail reads.
- Redirect behavior for deleted links is deferred until Stage 4.

Errors:

- `INVALID_ARGUMENT` for invalid ID.
- `NOT_FOUND` for missing or already soft-deleted links.
- `UNAUTHENTICATED` for missing or invalid bearer token.

## Stage 3 RPC Methods

Stage 3 RPC methods support administrator authentication and short-link management only. API handlers own HTTP parsing, response envelopes, and JWT creation or validation. RPC logic owns credential verification, data validation shared with management operations, and MySQL-backed mutations.

### AuthenticateAdmin

Authenticates administrator credentials and returns token subject data for API token creation.

Input:

- `username`
- `password`

Output:

- administrator ID
- username
- status

Errors:

- invalid credentials
- inactive administrator

### GetAdminProfile

Fetches the authenticated administrator profile by administrator ID.

### CreateShortLink

Creates a link, validates URL and code, generates a 6-character base62 code when code is omitted, and returns the created link.

### ListShortLinks

Returns paginated, non-deleted short-link summaries for management APIs.

### GetShortLink

Fetches non-deleted link details by ID.

### UpdateShortLink

Updates mutable link fields. Cache invalidation behavior is deferred until Stage 4.

### DeleteShortLink

Soft deletes a short link by setting `deleted_at`.

## Deferred Management APIs

The following management APIs are deferred until later stages.

### GET /admin/links/{id}/stats

Query parameters:

- `from`
- `to`

Returns totals and daily trend data.

## Deferred Redirect API

Redirect serving is deferred until Stage 4. Do not add this route during Stage 3.

### GET /{code}

Future results:

- `302 Found` with `Location` for active links.
- `404 Not Found` for missing or deleted links.
- `403 Forbidden` for disabled links.
- `410 Gone` for expired links.

## Deferred RPC Methods

### ResolveShortLink

Resolves a code into redirect data and domain status.

### RecordVisit

Records a visit event. Failure must be observable but should not block redirects.

### GetLinkStats

Returns aggregate statistics for a link and time range.

## Health Checks

### GET /healthz

Returns process liveness. This endpoint should not fail only because MySQL or Redis is temporarily unavailable.

### GET /readyz

Returns readiness. It checks dependencies needed to serve traffic, including RPC, MySQL, and Redis.

## Generation Boundary

Stage 3 management API and RPC contracts have generated go-zero skeleton code. Future contract changes
must update the `.api` or `.proto` source first, then run the approved `goctl` generation commands from
the repository root. Do not handwrite generated service code.
