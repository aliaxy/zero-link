# Testing Strategy

## Goals

- Treat tests as executable specifications.
- Keep unit tests fast and deterministic.
- Separate integration tests from normal unit tests.
- Verify the real redirect and cache behavior against MySQL and Redis.

## Unit Tests

Unit tests do not require external services.

Required areas:

- Short-code generation.
- Custom short-code validation.
- Reserved-word validation.
- URL validation.
- Expiration and status decisions.
- Error-to-HTTP mapping.
- Daily statistics aggregation.

Use table-driven tests with named subtests.

## Integration Tests

Integration tests use real MySQL and Redis and should be guarded by a build tag such as:

```go
//go:build integration
```

Required scenarios:

- Create a short link and read it from MySQL.
- Resolve a cache miss from MySQL and backfill Redis.
- Resolve a cache hit from Redis.
- Update a link and invalidate cache.
- Disable a link and confirm redirect is blocked.
- Expire a link and confirm `410 Gone`.
- Record a visit event.

## End-To-End Tests

E2E tests verify the complete local flow:

1. Start infrastructure.
2. Run migrations.
3. Start RPC.
4. Start API.
5. Log in through management API or UI.
6. Create a short link.
7. Visit the short URL.
8. Confirm redirect.
9. Confirm statistics appear.

## Race And Leak Detection

- Run `go test -race ./...` in CI once the module exists.
- Use goroutine leak detection if analytics workers or background queues are introduced.
- Background workers must accept context cancellation and shut down cleanly.

## Test Data

- Tests should generate unique short codes.
- Integration tests should clean their own data.
- Avoid relying on execution order.
- Use fixed clocks or injectable time providers for expiration tests.

## Acceptance Criteria

Before implementation is considered ready:

- Unit tests pass.
- Integration tests pass with local Compose dependencies.
- Redirect scenarios cover success, not found, disabled, expired, and cache miss.
- Management API scenarios cover authenticated and unauthenticated access.
