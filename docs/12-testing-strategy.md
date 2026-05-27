# Testing Strategy

## Goals

- Treat tests as executable specifications.
- Keep unit tests fast and deterministic.
- Separate integration tests from normal unit tests.
- Verify Stage 3 management behavior against MySQL before adding redirect and cache behavior.
- Keep redirect, Redis redirect cache, analytics, and management UI tests deferred until their implementation stages.

## Unit Tests

Unit tests do not require external services.

Required Stage 3 areas:

- Password hash verification accepts the seeded bcrypt hash format and rejects invalid passwords.
- Administrator authentication rejects missing credentials, invalid credentials, and inactive administrators.
- JWT creation includes administrator ID, username, issued time, and expiration.
- JWT validation rejects missing, malformed, expired, and incorrectly signed tokens.
- Short-code generation returns 6-character base62 values from cryptographically secure randomness.
- Automatic short-code creation retries on generated-code conflicts.
- Custom short-code validation accepts 3-32 characters containing ASCII letters, digits, `_`, and `-`.
- Custom short-code validation rejects empty values, short values, long values, spaces, slashes, non-ASCII characters, and reserved words.
- URL validation accepts absolute `http` and `https` URLs and rejects empty, relative, unsupported-scheme, and malformed URLs.
- Expiration validation accepts omitted or future timestamps and rejects past timestamps.
- Status validation accepts the Stage 3 active and disabled values.
- Error-to-HTTP mapping returns stable response envelope codes.
- Management middleware rejects unauthenticated requests before business logic runs.

Deferred unit areas:

- Redirect cache decisions.
- Redirect status decisions for missing, disabled, expired, or deleted links.
- Visit event recording.
- Daily statistics aggregation.
- Admin UI behavior.

Use table-driven tests with named subtests.

## Integration Tests

Integration tests use real MySQL and Redis and should be guarded by a build tag such as:

```go
//go:build integration
```

Required Stage 3 scenarios:

- Apply migrations with `golang-migrate`.
- Confirm the seeded local administrator exists for smoke testing.
- Log in successfully with an active administrator.
- Reject login with an invalid password.
- Reject login for inactive administrators.
- Reject management API requests without a Bearer token.
- Reject management API requests with an invalid or expired Bearer token.
- Return the authenticated administrator profile.
- Create a short link with a generated 6-character base62 code.
- Create a short link with a valid custom code.
- Reject invalid origin URLs.
- Reject reserved short codes.
- Reject duplicate custom short codes with `CONFLICT`.
- List short links with default pagination.
- List short links with explicit `page`, `page_size`, `status`, and `keyword` filters.
- Fetch short-link details by ID.
- Update mutable short-link fields while keeping the code immutable.
- Reject update requests that attempt to change the code.
- Soft delete a short link.
- Hide soft-deleted links from normal management list and detail reads.

Deferred integration scenarios:

- Resolve a cache miss from MySQL and backfill Redis.
- Resolve a cache hit from Redis.
- Disable a link and confirm redirect is blocked.
- Expire a link and confirm `410 Gone`.
- Record a visit event.
- Query link statistics.
- Exercise management UI workflows.

## End-To-End Smoke Tests

Stage 3 smoke tests verify the complete local management flow:

1. Start infrastructure.
2. Run migrations.
3. Start `link-rpc`.
4. Start `link-api`.
5. Log in through `POST /admin/login` with the seeded local administrator.
6. Call `GET /admin/profile` with the returned bearer token.
7. Create a short link through `POST /admin/links`.
8. List short links through `GET /admin/links`.
9. Fetch the created short link through `GET /admin/links/{id}`.
10. Update mutable fields through `PATCH /admin/links/{id}`.
11. Soft delete through `DELETE /admin/links/{id}`.
12. Confirm the deleted link is hidden from normal management reads.

Deferred E2E steps:

- Visit the short URL.
- Confirm redirect.
- Confirm cache hit and miss behavior.
- Confirm statistics appear.
- Complete the workflow through the management UI.

## Documentation Review Checks

Before Stage 3 implementation starts, the documentation-only specification pass is accepted when:

- Every Stage 3 HTTP endpoint documents authentication, request shape, response shape, and error behavior.
- Every Stage 3 RPC method has a matching management use case.
- JWT configuration is documented without committing real secrets.
- Pagination defaults and maximum page size are documented.
- Short-code generation and custom-code validation rules are documented.
- Duplicate custom codes map to `CONFLICT`.
- Redirect serving, Redis redirect cache, analytics, and management UI remain explicitly out of scope.
- The implementation handoff says `goctl` generation must be explicitly run in the next implementation pass after specification review.

## Race And Leak Detection

- Run `go test -race ./...` in CI once the Stage 3 implementation is in place.
- Use goroutine leak detection if analytics workers or background queues are introduced.
- Background workers must accept context cancellation and shut down cleanly.

## Test Data

- Tests should generate unique short codes unless they intentionally verify conflict behavior.
- Integration tests should clean their own data.
- Tests must not depend on execution order.
- Use fixed clocks or injectable time providers for expiration and token tests.
- Do not depend on the local/dev seeded administrator outside integration and smoke-test flows.
- Do not log plaintext passwords, bearer tokens, or real secrets in test output.

## Acceptance Criteria

Before Stage 3 implementation is considered ready:

- Unit tests pass.
- Integration tests pass with local Compose dependencies.
- Stage 3 management scenarios cover authenticated and unauthenticated access.
- Stage 3 administrator scenarios cover login, invalid credentials, inactive administrators, JWT creation, JWT validation, and profile retrieval.
- Stage 3 short-link scenarios cover creation validation, generated codes, custom codes, duplicate code conflict, pagination, detail, update, immutable code, and soft delete.
- Redirect, cache, analytics, and UI scenarios remain deferred until their dedicated stages.
