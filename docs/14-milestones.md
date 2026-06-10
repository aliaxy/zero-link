# Milestones

## Stage 0: Requirements And Technical Documents

Goal: establish the product and technical baseline.

Inputs:

- Project intent.
- go-zero, MySQL, Redis, Docker Compose constraints.

Outputs:

- Documentation set under `docs/`.
- Clear scope for implementation planning.

Acceptance:

- Product scope, architecture, storage, API, analytics, security, testing, deployment, and roadmap are documented.

Risks:

- Over-designing before validating the core flow.
- Leaving ambiguous ownership between API and RPC services.

## Stage 1: Engineering Initialization And Local Dependencies

Goal: create a skeleton-only development foundation.

Inputs:

- Documentation baseline.

Outputs:

- Go module.
- Reserved service, migration, config, and admin UI directories.
- `deploy/docker-compose.infra.yml`.
- MySQL and Redis local setup.
- Example and local configuration files.
- Placeholder documentation for future generated service contracts and migrations.

Acceptance:

- `make test` reports that the skeleton foundation is ready.
- `deploy/docker-compose.infra.yml` validates with Docker Compose.
- MySQL and Redis can be started locally with `make infra-up`.
- No generated go-zero service code, API route contracts, protobuf RPC methods, business migrations, or short-link logic exist yet.

Risks:

- Local toolchain not loaded outside the Nix shell.
- Introducing too many services too early.

## Stage 2: API/RPC Skeleton

Goal: lock transport boundaries.

Outputs:

- API route definitions.
- RPC proto definitions.
- Generated go-zero API and RPC service skeletons.
- Basic health checks.
- Shared error response mapping.
- Minimal local API-to-RPC wiring.

Acceptance:

- `healthz` and `readyz` work.
- API-to-RPC call succeeds locally.

Risks:

- Transport types leaking into domain logic.

## Stage 3: Administrator And Link Management

Goal: create and manage short links.

Status: backend management API/RPC implementation and local HTTP smoke request assets are in place.

Outputs:

- Administrator login.
- Create link.
- List links.
- View link details.
- Update mutable fields.
- Soft delete.
- Local httpyac smoke requests for health and management flows.

Acceptance:

- Authenticated administrator can manage links.
- Unauthenticated management requests fail.
- Duplicate custom codes return conflict.

Risks:

- Weak password storage.
- Missing cache invalidation after updates.

## Stage 4: Redirect And Cache

Goal: deliver fast and correct redirects.

Status: complete.

Outputs:

- `GET /{code}` route.
- Redis cache lookup and backfill.
- Not found, disabled, and expired handling.
- Cache invalidation on link updates.

Acceptance:

- Active links redirect.
- Missing links return 404.
- Disabled links return 403.
- Expired links return 410.
- Cache miss and cache hit paths both pass tests.

Risks:

- Stale cache after updates.
- Synchronous work slowing redirects.

## Stage 5: Analytics

Goal: record and display basic visit statistics.

Status: complete.

Outputs:

- Visit event recording via `AnalyticsMiddleware` (non-blocking, goroutine per redirect).
- Daily stat aggregation with `UpsertPV` (ON DUPLICATE KEY UPDATE).
- `GET /admin/links/{id}/stats` statistics query API.

Acceptance:

- Visits are recorded after redirects.
- Statistics are queryable via the management API.
- Analytics failures do not block redirects.

Risks:

- Event table grows quickly.
- UV definition changes later (Stage 5 approximation: uv stays at 1 after first insert).

## Stage 6: Management UI

Goal: provide an administrator-facing workflow.

Status: complete.

Outputs:

- Login page.
- Link list.
- Create/edit form.
- Link details and stats view.

Acceptance:

- Administrator can complete the core workflow without direct API calls.

Risks:

- Spending too much time on UI polish before backend behavior is stable.

## Stage 7: Testing, Observability, And Security Hardening

Goal: make the service diagnosable and safer.

Status: complete.

Completed:

- Reserved short codes extended with `metrics` and `static`.
- `MaxConns` and `MaxBytes` connection protection on `link-api`.
- IP rate limiting on redirect (20 req/s) and login (10 req/min) via go-zero `PeriodLimit`.
- Structured logging across `link-api` (login, auth, redirect) and `link-rpc` (resolve, authenticate, analytics).
- Prometheus business metrics: `zerolink_redirect_requests_total{result}` and `zerolink_analytics_events_total{result}`.
- Prometheus `/metrics` endpoints on ports 9100 (link-api) and 9101 (link-rpc).
- Unit tests for `AnalyticsMiddleware`, `IPRateLimitMiddleware`, `LoginRateLimitMiddleware`.
- Full-stack integration tests in `tests/integration/` run via `make test-e2e`.

Outputs:

- Unit and integration tests.
- Structured logs.
- Metrics.
- Health checks.
- Rate limits.
- Security review fixes.

Acceptance:

- Unit tests pass.
- Integration tests pass against local dependencies.
- Common failures are visible in logs and metrics.

Risks:

- Metrics label cardinality mistakes.
- Duplicate error logging.

## Stage 8: Cache Optimization And Data Lifecycle

Goal: defend against cache penetration and enforce data retention.

Status: complete.

Completed:

- In-process cuckoo filter in `services/link-rpc/pkg/filter` loaded from `short_link` on startup (including soft-deleted rows).
- Redis Pub/Sub channel `zl:code:created` synchronises filter updates across multiple `link-rpc` instances; subscribe before batch load to close the race window.
- `ResolveShortLink` fast path: filter miss returns `NotFound` immediately without touching Redis or MySQL.
- `CreateShortLink` inserts new codes into the local filter and publishes to the Pub/Sub channel after a successful create.
- Cache breakdown already covered by go-zero's internal `syncx.SingleFlight` inside `QueryRowIndexCtx`.
- Migration `000003_data_retention`: `short_link_archive` table (permanent archive with original `id`) and `reserved_code` table (prevents code reuse; primary key on `code`).
- `reserved_code` model generated by goctl (`--style gozero`); custom `Exists` method added in the non-generated file.
- `CreateShortLink` checks `reserved_code` for custom codes to prevent reuse of archived codes.
- Cleanup runner (`internal/cleanup/`) runs on startup and every 24 hours: archives soft-deleted `short_link` rows after 365 days, batch-deletes `visit_event` after 90 days, batch-deletes `link_daily_stat` after 730 days.
- Archival step is idempotent: `INSERT IGNORE` at each step (no transaction needed).
- Prometheus counters: `zerolink_filter_requests_total{result}` and `zerolink_cleanup_deleted_rows_total{table}`.
- Unit tests: cuckoo filter (basic and concurrent), cleanup runner cancellation (goleak), reserved code conflict, filter miss/hit in resolve logic.

Outputs:

- `migrations/000003_data_retention.up.sql` / `.down.sql`.
- `services/link-rpc/pkg/filter/`.
- `services/link-rpc/internal/cleanup/`.
- `services/link-rpc/internal/model/reservedcodemodel.go` (goctl-generated + custom `Exists`).

Acceptance:

- Filter blocks unknown codes without touching Redis or DB.
- Newly created codes immediately resolve correctly after filter update.
- Archived links are no longer in `short_link`; their codes remain in `reserved_code`.
- Cleanup metrics visible in Prometheus.

Risks:

- Cuckoo filter false positives allow a small number of non-existent codes through to Redis/MySQL (acceptable; these fall through to the existing not-found path).
- Filter is not persisted; a restart reloads from DB (acceptable given startup time).

## Stage 9: Full Docker Compose Deployment

Goal: package the full system for deployment.

Status: complete.

Outputs:

- Dockerfiles for API and RPC.
- Complete `docker-compose.yml`.
- Deployment documentation.

Acceptance:

- Full Compose stack starts.
- Data persists across restarts.
- Core E2E flow works in Compose.

Risks:

- Config drift between local and Compose modes.
- Service readiness confused with container startup order.

## Post-Stage-9: Production Hardening

These improvements were made after Stage 9 without advancing the feature stage boundary. Each maps to a
named sub-stage and corresponds to its own feature branch merged into main.

### Stage A: Atomic Archive And Cache Invalidation

Goal: fix two P0 data-correctness issues in the cleanup and cache layers.

Status: complete.

Completed:

- `cleanup/shortlinks.go:archiveLink` wrapped in `db.TransactCtx`: INSERT archive → INSERT reserved_code →
  DELETE short_link now succeed or fail atomically. A process crash can no longer leave an un-reserved code.
- After a successful transaction, the cleanup runner calls `rdb.Del` on the id and code cache keys so that
  archived links are not served stale from Redis.
- Cache key constants (`CacheShortLinkIdPrefix`, `CacheShortLinkCodePrefix`) exported from
  `model/vars.go` and reused by the cleanup package.
- `cleanup.NewRunner` now accepts a `*redis.Redis` parameter; `svc/servicecontext.go` passes `rdb`.
- `model/linkarchiver.go` introduced to keep the archival transaction logic separate from the goctl-generated
  model files.

### Stage B: Dual-Token Authentication And Password Change

Goal: fix JWT revocation gap and add administrator password management.

Status: complete.

Completed:

- Dual-token architecture: stateless access token (JWT, 15 min) paired with an opaque refresh token
  (7 days, Redis-backed). Access token statelessness is an accepted industry trade-off; short TTL limits
  the exposure window.
- `services/link-api/internal/auth/refreshtoken.go` — `RefreshTokenStore` implements `RefreshTokenIssuer`:
  - `Issue`: generates a 32-byte `crypto/rand` token, stores the SHA-256 hash in `zl:rt:{hash}` (TTL 7d)
    and adds the hash to the per-user set `zl:rt:user:{adminID}` (TTL 7d).
  - `Rotate`: verifies the hash exists, deletes it, issues a new token (reuse detection).
  - `RevokeAll`: reads `SMEMBERS zl:rt:user:{adminID}`, bulk-deletes all token keys and the user set.
- `POST /admin/login` now returns `access_token`, `access_token_expires_at`, `refresh_token`,
  `refresh_token_expires_at`, and `admin`.
- `POST /admin/refresh` — rotates the refresh token and issues a new access token; invalid or expired
  tokens return `UNAUTHENTICATED 401`.
- `PATCH /admin/password` — authenticated endpoint; calls `ChangePassword` RPC, then `RevokeAll` to
  invalidate all sessions.
- `ChangePassword` RPC added to `AdminService` in `link.proto`; verifies old password, hashes new
  password with bcrypt, updates `admin_user.password_hash`.
- httpyac smoke requests updated: `token` → `accessToken`, `refreshToken` added, refresh and
  change-password requests added.

### Stage C: Performance Critical Path

Goal: eliminate three structural performance problems in analytics, UV queries, and filter loading.

Status: complete.

Completed:

- `AnalyticsMiddleware` replaced unbounded goroutine-per-redirect with a channel-backed worker pool
  (8 workers, 2 000-slot buffer). Events are dropped gracefully when the buffer is full; `Stop()` closes
  the channel for clean shutdown.
- `HasVisitedToday` rewritten from `date(visited_at) = ?` to a half-open range
  `visited_at >= dayStart AND visited_at < dayEnd`. MySQL can now use the `visited_at` index instead of
  evaluating `date()` on every row.
- `loadCodesIntoFilter` switched from `LIMIT n OFFSET m` to id-cursor pagination
  (`WHERE id > ? ORDER BY id LIMIT ?`). Eliminates O(n) offset scans on large tables.
