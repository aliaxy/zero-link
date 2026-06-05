# Error Handling

## Principles

- Errors are either returned or logged, not both.
- Internal errors should be wrapped with context.
- HTTP and RPC boundaries translate internal errors into stable public error codes.
- Technical details are not exposed to administrators or visitors.
- Expected business failures are not panics.

## Error Categories

Management API JSON error codes (the `code` field in the response envelope):

- `BAD_REQUEST`: request payload could not be parsed (missing required field, type mismatch, malformed JSON). Produced by go-zero's `httpx.Parse` before the RPC call is made.
- `INVALID_ARGUMENT`: request passed parsing but contains an invalid value (e.g. malformed URL, invalid short code).
- `UNAUTHENTICATED`: missing or invalid bearer token, or RPC authentication failure.
- `PERMISSION_DENIED`: authenticated but not allowed, or disabled public link.
- `NOT_FOUND`: resource does not exist.
- `CONFLICT`: unique constraint or state conflict.
- `GONE`: short link has expired or been disabled.
- `SERVICE_UNAVAILABLE`: MySQL, Redis, or RPC is not reachable.
- `INTERNAL`: unexpected server error.

## HTTP Mapping

Management APIs:

- `BAD_REQUEST`: 400.
- `INVALID_ARGUMENT`: 400.
- `UNAUTHENTICATED`: 401.
- `PERMISSION_DENIED`: 403.
- `NOT_FOUND`: 404.
- `CONFLICT`: 409.
- `GONE`: 410.
- `SERVICE_UNAVAILABLE`: 503.
- `INTERNAL`: 500.

Redirect API (plain text, no JSON envelope):

- Missing or deleted link: 404.
- Disabled link: 403.
- Expired link: 410.
- Unexpected failure: 500.

## Parse Error Handling

go-zero's `httpx.Parse` runs before logic and produces plain Go errors (not gRPC status errors) when a required JSON field is absent, the value has the wrong type, or the body is malformed. The management API error handler (`apierror.ErrorHandler`) catches these patterns and returns `BAD_REQUEST 400` so the caller is not presented with a generic `INTERNAL 500`.

## Go Error Handling

Implementation should use sentinel or typed domain errors for expected business cases.

Examples:

- `ErrShortLinkNotFound`
- `ErrShortLinkDisabled`
- `ErrShortLinkExpired`
- `ErrShortCodeConflict`
- `ErrInvalidOriginURL`

Lower layers wrap errors with context:

```go
return fmt.Errorf("find short link by code: %w", err)
```

Boundary layers inspect with `errors.Is` or `errors.As` and convert to response codes.

## Logging Boundary

Recommended logging points:

- HTTP middleware for request summary.
- HTTP error response boundary.
- RPC interceptor or method boundary.
- Background worker failure boundary.

Do not log the same error in every helper function.

## User-Facing Messages

Messages should be stable and understandable:

- `short link not found`
- `short link is disabled`
- `short link has expired`
- `invalid origin url`
- `short code already exists`

Do not return SQL errors, Redis errors, stack traces, or dependency addresses in public responses.
