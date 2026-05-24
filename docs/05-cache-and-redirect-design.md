# Cache And Redirect Design

## Redis Key Design

- `shortlink:code:{code}`: cached short-link resolution payload.
- `ratelimit:create:{admin_id}`: administrator create-rate counter.
- `ratelimit:redirect:{code}:{ip_hash}`: per-link redirect-rate counter.
- `uv:{link_id}:{date}`: optional UV de-duplication set or bitmap.

## Cached Payload

`shortlink:code:{code}` should contain only fields needed by redirect resolution:

- `link_id`
- `code`
- `origin_url`
- `status`
- `expire_at`

The payload must not contain administrator secrets or unrelated metadata.

## Redirect Flow

1. Receive `GET /{code}`.
2. Validate short-code format.
3. Apply lightweight redirect rate limiting.
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

Delete `shortlink:code:{code}` after:

- Original URL update.
- Status update.
- Expiration time update.
- Soft delete.

Cache invalidation should be best effort but observable. If Redis deletion fails, the update operation should return an error unless product requirements explicitly allow temporary stale redirects.

## Performance Boundary

The redirect path must do the minimum work needed to decide the redirect result. Heavy analytics parsing, geographic lookup, and aggregation run outside the synchronous redirect path.
