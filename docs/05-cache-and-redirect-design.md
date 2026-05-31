# Cache And Redirect Design

## Redis Key Design

- `cache:shortLink:id:{id}`: short-link row cached by primary key (goctl model layer).
- `cache:shortLink:code:{code}`: short-link row cached by unique code index (goctl model layer).
- `ratelimit:create:{admin_id}`: administrator create-rate counter.
- `rl:redirect:ip:{ip}`: per-IP redirect rate counter (go-zero `PeriodLimit`, 20 req/s window).
- `rl:login:ip:{ip}`: per-IP login rate counter (go-zero `PeriodLimit`, 10 req/min window).
- `uv:{link_id}:{date}`: optional UV de-duplication set or bitmap.

The `cache:shortLink:*` keys are managed by the go-zero goctl cached model. `FindOneByCode` uses a
two-level index cache (code → id → full row). `Update` and `Delete` invalidate both keys automatically.
Custom operations that bypass the generated `Update` or `Delete` (such as soft delete) must invalidate
both keys explicitly.

## Cached Payload

`cache:shortLink:code:{code}` caches the full `short_link` row. The redirect logic reads `origin_url`,
`status`, and `expire_at` from the cached row. Fields unrelated to redirect resolution are present but
unused on the redirect path.

## Redirect Flow

1. Receive `GET /{code}`.
2. Validate short-code format.
3. Apply per-IP rate limiting (20 req/s; returns `429 Too Many Requests` on excess).
4. Read `shortlink:code:{code}` from Redis.
5. If Redis misses, read `short_link` from MySQL by `code`.
6. If no row exists, return `404 Not Found`.
7. If the link is disabled, return `403 Forbidden`.
8. If the link is expired, return `410 Gone`.
9. If the link is active, backfill Redis with a safe TTL.
10. Emit a visit event asynchronously.
11. Return `302 Found` with `Location: {origin_url}`.

## Error Mapping

- Invalid code format: `404 Not Found`.
- Missing link: `404 Not Found`.
- Deleted link: `404 Not Found`.
- Disabled link: `403 Forbidden`.
- Expired link: `410 Gone`.
- Internal dependency failure: `500 Internal Server Error` for API clients, generic error page for browser visitors.

## Cache Expiration

- Links without business expiration can use a default cache TTL such as 10 minutes.
- Links with `expire_at` must use a TTL no longer than the remaining lifetime.
- Disabled and deleted links should not be cached as successful redirect payloads in the first stage.
- Negative caching can be added later if missing-link traffic becomes expensive.

## Cache Invalidation

Invalidate `cache:shortLink:id:{id}` and `cache:shortLink:code:{code}` after:

- Original URL update.
- Status update.
- Expiration time update.
- Soft delete.

The goctl `Update` and `Delete` methods handle invalidation automatically. Soft delete uses a custom
`ExecCtx` call that invalidates both keys explicitly before executing the SQL. Cache invalidation is
best effort; if Redis deletion fails the operation still returns success to the caller.

## Performance Boundary

The redirect path must do the minimum work needed to decide the redirect result. Heavy analytics parsing, geographic lookup, and aggregation run outside the synchronous redirect path.
