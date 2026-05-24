# Error Handling

## Principles

- Errors are either returned or logged, not both.
- Internal errors should be wrapped with context.
- HTTP and RPC boundaries translate internal errors into stable public error codes.
- Technical details are not exposed to administrators or visitors.
- Expected business failures are not panics.

## Error Categories

- `INVALID_ARGUMENT`: invalid input such as malformed URL or code.
- `UNAUTHORIZED`: missing or invalid authentication.
- `FORBIDDEN`: authenticated but not allowed, or disabled public link.
- `NOT_FOUND`: resource does not exist.
- `CONFLICT`: unique constraint or state conflict.
- `EXPIRED`: short link has expired.
- `DISABLED`: short link has been disabled.
- `DEPENDENCY_UNAVAILABLE`: MySQL, Redis, or RPC unavailable.
- `INTERNAL`: unexpected server error.

## HTTP Mapping

Management APIs:

- `INVALID_ARGUMENT`: 400.
- `UNAUTHORIZED`: 401.
- `FORBIDDEN`: 403.
- `NOT_FOUND`: 404.
- `CONFLICT`: 409.
- `DEPENDENCY_UNAVAILABLE`: 503.
- `INTERNAL`: 500.

Redirect API:

- Missing or deleted link: 404.
- Disabled link: 403.
- Expired link: 410.
- Unexpected failure: 500.

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
