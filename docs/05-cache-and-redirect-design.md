# Cache And Redirect Design

## Redis Key Design

- `cache:shortLink:id:{id}`: short-link row cached by primary key (goctl model layer).
- `cache:shortLink:code:{code}`: short-link row cached by unique code index (goctl model layer).
- `ratelimit:create:{admin_id}`: administrator create-rate counter.
- `rl:redirect:ip:{ip}`: per-IP redirect rate counter (go-zero `PeriodLimit`, 20 req/s window).
- `rl:login:ip:{ip}`: per-IP login rate counter (go-zero `PeriodLimit`, 10 req/min window).
- `uv:{link_id}:{date}`: optional UV de-duplication set or bitmap.
- `zl:code:created` (Pub/Sub channel): published by `link-rpc` after every successful short-link creation; subscribers update their local cuckoo filter.

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

## Cache Penetration Defence

Cache penetration occurs when repeated requests for non-existent codes hit Redis and MySQL on every request. `ResolveShortLink` defends against this with an in-process cuckoo filter.

### Cuckoo Filter

The filter lives in `services/link-rpc/pkg/filter` and is loaded at startup:

1. Subscribe to `zl:code:created` before the batch read to close the race window.
2. Batch-read all `code` values from `short_link` (including soft-deleted rows, which still occupy the unique index).
3. Insert each code into the filter.

On `ResolveShortLink`, the filter is checked before any Redis or MySQL access:

- **Filter miss** (definitive non-existence): return `NotFound` immediately without touching Redis or DB. Increments `zerolink_filter_requests_total{result="miss"}`.
- **Filter hit** (probable existence): proceed to the normal Redis → MySQL path. Increments `zerolink_filter_requests_total{result="hit"}`.

Cuckoo filters have a small false-positive rate (a hit does not guarantee existence) but zero false negatives, making them safe for this use case.

### Multi-Instance Synchronisation

Multiple `link-rpc` instances each hold an independent in-process filter. When `CreateShortLink` succeeds, the creating instance:

1. Inserts the new code into its own filter.
2. Publishes the code to the `zl:code:created` Redis Pub/Sub channel.

Each instance runs a subscription goroutine that inserts published codes into its local filter. A failed publish is non-fatal: the other instances' filters miss the new code temporarily, and those requests fall through to the normal Redis → MySQL path correctly.

### Cache Breakdown

Cache breakdown (many concurrent requests for the same valid code on a cache miss) is already handled by go-zero's internal `syncx.SingleFlight` inside `QueryRowIndexCtx`. No additional fix is needed.
