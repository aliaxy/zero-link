# Testing Strategy

## Goals

- Treat tests as executable specifications.
- Keep unit tests fast and deterministic.
- Separate integration tests from normal unit tests.
- Keep management UI tests deferred until Stage 6.

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

Implemented unit areas (Stage 4 and Stage 5):

- Redirect cache decisions.
- Redirect status decisions for missing, disabled, expired, or deleted links.
- Visit event recording and daily statistics aggregation (analytics helpers, RecordVisit, GetLinkStats).
- AnalyticsMiddleware goroutine triggering and IP extraction.

Implemented unit areas (Stage 7):

- `IPRateLimitMiddleware`: allowed, over-quota, limiter error, hit-quota pass-through, XFF key extraction, RemoteAddr fallback. Uses `rateLimiter` interface stub to avoid real Redis.
- `LoginRateLimitMiddleware`: allowed and over-quota via `IPRateLimitMiddleware` delegation.
- `AnalyticsMiddleware`: 302 triggers `RecordVisit`, 404/403 do not. Uses `goleak.VerifyTestMain` to detect goroutine leaks.
- `ExtractClientIP`: XFF header (first IP from comma-separated list), RemoteAddr fallback.

Implemented unit areas (Stage 8):

- `pkg/filter/cuckoo_test.go`: filter returns false before insert, true after insert, and passes `-race` on concurrent insert and lookup.
- `cleanup/runner_test.go`: `TestRunner_StartCancellation` verifies the background goroutine exits cleanly when context is cancelled. Uses `goleak.VerifyTestMain(m, goleak.IgnoreCurrent())` to detect leaks against go-zero framework goroutines.
- `linkservice/createshortlinklogic_test.go`: `TestCreateShortLink_ReservedCode_Conflict` verifies that a code absent from `short_link` but present in `reserved_code` returns `AlreadyExists`.
- `linkservice/resolveshortlinklogic_test.go`: `TestResolveShortLink_FilterMiss_SkipsDB` verifies that an empty filter returns `NotFound` without reaching the model. `TestResolveShortLink_FilterHit_ProceedsToModel` verifies that a filter containing the code proceeds to the model and returns the origin URL.

Deferred unit areas:

- Admin UI behavior.

## Full-Stack Integration Tests

Full-stack integration tests live under `tests/integration/` and exercise the complete HTTP → link-api → link-rpc → MySQL/Redis chain. They are guarded by `//go:build integration` and run with:

```bash
make test-e2e
```

Prerequisites: `make infra-up && make migrate-up`.

`TestMain` starts both services as OS subprocesses (`go run ./services/link-rpc` and `go run ./services/link-api`) on dynamically allocated ports, waits for readiness, then runs all tests in the package. On teardown the subprocess tree is killed via process group (`Setpgid + SIGKILL`) to prevent pipe-close hangs.

Implemented scenarios:

- `GET /healthz` and `GET /readyz` return `{status:"ok"}`.
- `POST /admin/login` succeeds with the migration-seeded admin and returns a JWT.
- `POST /admin/login` returns 401 with `UNAUTHENTICATED` code on wrong password.
- `GET /admin/profile` returns the authenticated admin's username.
- All seven management endpoints return 401 without a token.
- Link full CRUD: create → list (link appears) → get (fields match) → update (new URL and title) → delete (deleted=true).
- `GET /admin/links/999999999` returns 404.
- `GET /{code}` on an active link returns 302 with correct `Location`.
- `GET /{code}` on a missing code returns 404.
- `GET /{code}` on a disabled link returns 403.
- `GET /{code}` on an expired link returns 410. Expiry is set via DB `UPDATE` because `ValidateExpireAt` rejects past timestamps through the API.
- `RecordVisit` + `GetLinkStats`: redirect fires the analytics goroutine; test polls `/admin/links/:id/stats` up to 8 s until `pv ≥ 1`.
- Multiple redirects: poll up to 10 s until `pv ≥ 3`.
- New link with no redirects returns an empty `items` array from stats.

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

- Exercise management UI workflows (deferred to a future stage).

## End-To-End Smoke Tests

Smoke tests verify the complete local flow:

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
13. Access `GET /{code}` with an active link and confirm `302`.
14. Confirm cache hit and miss behavior on repeated requests.
15. Call `GET /admin/links/{id}/stats` and confirm `pv: 1` for today.

Deferred E2E steps:

- Complete the management workflow through the UI.

## Documentation And Smoke Asset Checks

Documentation and smoke assets are accepted when:

- Every HTTP endpoint documents authentication, request shape, response shape, and error behavior.
- Every RPC method has a matching use case.
- JWT and IP salt configuration are documented without committing real secrets.
- Pagination defaults and maximum page size are documented.
- Short-code generation and custom-code validation rules are documented.
- Analytics date range defaults and maximum range are documented.
- HTTP smoke requests live under `tests/httpyac/`.
- `tests/httpyac/http-client.example.env.json` documents the `local` httpyac environment.
- `tests/httpyac/http-client.private.env.json` stores local passwords, bearer tokens, and generated IDs and must not be committed.

## Race And Leak Detection

- Run `go test -race ./...` before merging implementation changes.
- `AnalyticsMiddleware` tests use `go.uber.org/goleak` to verify no goroutines leak after the 3-second timeout expires.
- Background goroutines must have a clear exit condition (context cancellation or timeout).

## Test Data

- Tests should generate unique short codes unless they intentionally verify conflict behavior.
- Integration tests should clean their own data.
- Tests must not depend on execution order.
- Use fixed clocks or injectable time providers for expiration and token tests.
- Do not depend on the local/dev seeded administrator outside integration and smoke-test flows.
- Do not log plaintext passwords, bearer tokens, or real secrets in test output.

## Acceptance Criteria

Before implementation is considered ready:

- Unit tests pass.
- Integration tests pass with local Compose dependencies.
- Management scenarios cover authenticated and unauthenticated access.
- Administrator scenarios cover login, invalid credentials, inactive administrators, JWT creation, JWT validation, and profile retrieval.
- Short-link scenarios cover creation validation, generated codes, custom codes, duplicate code conflict, pagination, detail, update, immutable code, and soft delete.
- Redirect scenarios cover active, missing, disabled, and expired links.
- Analytics scenarios cover visit recording, PV upsert, stats query, date range defaults, and invalid range rejection.
- Management UI scenarios remain deferred to Stage 6.
