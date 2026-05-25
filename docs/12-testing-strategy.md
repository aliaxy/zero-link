# Testing Strategy

## Goals

- Treat tests as executable specifications.
- Keep unit tests fast and deterministic.
- Separate integration tests from normal unit tests.
- Verify Stage 3 management behavior against MySQL before adding redirect and cache behavior.

## Unit Tests

Unit tests do not require external services.

Required areas:

- Password hash verification.
- JWT creation and validation.
- Short-code generation.
- Custom short-code validation.
- Reserved-word validation.
- URL validation.
- Expiration and status decisions.
- Error-to-HTTP mapping.

Deferred unit areas:

- Redirect cache decisions.
- Daily statistics aggregation.

Use table-driven tests with named subtests.

## Integration Tests

Integration tests use real MySQL and Redis and should be guarded by a build tag such as:

```go
//go:build integration
```

Required scenarios:

- Apply migrations with `golang-migrate`.
- Seed or create an administrator and log in successfully.
- Reject login with an invalid password.
- Reject management API requests without a Bearer token.
- Create a short link and read it from MySQL.
- Reject invalid origin URLs and reserved short codes.
- Reject duplicate custom short codes with a conflict.
- List short links with pagination.
- Fetch short-link details.
- Update mutable short-link fields while keeping the code immutable.
- Soft delete a short link and hide it from normal management reads.

Deferred integration scenarios:

- Resolve a cache miss from MySQL and backfill Redis.
- Resolve a cache hit from Redis.
- Disable a link and confirm redirect is blocked.
- Expire a link and confirm `410 Gone`.
- Record a visit event.

## End-To-End Tests

E2E tests verify the complete local flow:

1. Start infrastructure.
2. Run migrations.
3. Start RPC.
4. Start API.
5. Log in through the management API with the seeded local administrator.
6. Create a short link.

Deferred E2E steps:

- Visit the short URL.
- Confirm redirect.
- Confirm statistics appear.

## Race And Leak Detection

- Run `go test -race ./...` in CI once the module exists.
- Use goroutine leak detection if analytics workers or background queues are introduced.
- Background workers must accept context cancellation and shut down cleanly.

## Test Data

- Tests should generate unique short codes.
- Integration tests should clean their own data.
- Avoid relying on execution order.
- Use fixed clocks or injectable time providers for expiration tests.
- Do not depend on the local/dev seeded administrator outside integration and smoke-test flows.

## Acceptance Criteria

Before implementation is considered ready:

- Unit tests pass.
- Integration tests pass with local Compose dependencies.
- Stage 3 management scenarios cover authenticated and unauthenticated access.
- Stage 3 short-link scenarios cover creation validation, code conflict, pagination, detail, update, and soft delete.
- Redirect scenarios cover success, not found, disabled, expired, and cache miss before Stage 4 completion.
