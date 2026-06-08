# Implementation Plan

## Goal

Track the current zero-link implementation state and define the next safe handoff point. Stages 1–9 are
complete.

## Source Documents

This plan is derived from:

- `docs/00-project-overview.md`
- `docs/03-system-architecture.md`
- `docs/04-storage-design.md`
- `docs/06-api-design.md`
- `docs/10-local-development.md`
- `docs/12-testing-strategy.md`
- `docs/14-milestones.md`
- `docs/17-git-workflow.md`
- `AGENTS.md`

## Current Status

Stage 6 management UI is complete. The implementation includes:

- Vue 3 + Vite + Element Plus admin SPA under `web/admin/`.
- Apple HIG aesthetic: frosted glass sidebar, system blue accent, Geist Variable font.
- Login page with JWT stored in localStorage; 401 auto-redirect via axios interceptor.
- Links list with keyword search, status filter, row-click to detail, and action buttons.
- Link create/edit via `LinkFormDrawer` with client-side validation (URL format, code pattern, reserved words).
- Link detail view with copy short URL and navigate to stats.
- Stats view with ECharts dual-axis PV/UV line chart, summary cards, date range picker, and data table.
- CORS middleware in `link-api` applied globally via `server.Use()` with configurable `AllowOrigins`.
- Vite dev proxy: `/api/*` strips prefix and forwards to `link-api`; `/:code` regex proxy forwards short link codes to `link-api` for redirect.
- Makefile targets: `web-install`, `web-dev`, `web-build`, `web-preview`.
- `web/admin/.env.example` documents `VITE_API_BASE_URL`, `VITE_API_TARGET`, `VITE_SHORT_LINK_BASE`.

Stage 1 is complete. The repository foundation includes:

- Go module `github.com/aliaxy/zero-link`.
- Local dependency Compose file at `deploy/docker-compose.infra.yml`.
- Example configuration files under `etc/`.
- Reserved directories for services and future admin UI.
- Stage 3 foundation migrations under `migrations/`.
- Make targets for local infrastructure, service runs, tests, formatting, and linting.

Stage 2 is complete. The transport skeleton includes:

- go-zero API service under `services/link-api/`.
- go-zero RPC service under `services/link-rpc/`.
- Health endpoints `GET /healthz` and `GET /readyz`.
- RPC readiness method `HealthService.Check`.
- API readiness logic that calls the RPC readiness method.
- RPC readiness logic that validates configured MySQL and Redis endpoints with simple connectivity checks.

Stage 5 analytics is complete. The implementation includes:

- `visit_event` and `link_daily_stat` MySQL tables under `migrations/000002_stage5_analytics`.
- `VisitEventModel` and `LinkDailyStatModel` (goctl-generated plain sqlx models, no Redis cache).
- `RecordVisit` RPC: IP hashing (HMAC-SHA256 + salt), device detection, visit event insert, daily PV upsert.
- `GetLinkStats` RPC: date range validation (default 30 days, max 90 days), daily stat query.
- `AnalyticsMiddleware` in `link-api`: wraps redirect route via `statusRecorder`, fires goroutine on 302 only with 3s timeout context; non-302 responses are ignored.
- `GET /admin/links/{id}/stats` endpoint returning daily PV/UV in standard envelope.
- `Analytics.IPSalt` config in `link-rpc` only; never exposed to `link-api`.

Stage 4 redirect and cache is complete. The implementation includes:

- `ResolveShortLink` RPC with goctl two-level Redis cached model (`FindOneByCode`).
- `CacheRedis` config wired to `AdminUserModel` and `ShortLinkModel`.
- `ErrPermissionDenied` and `ErrGone` domain errors mapped to gRPC `PermissionDenied` and `FailedPrecondition`.
- Cache invalidation on `Update` and `SoftDelete` via goctl `ExecCtx` with id and code keys.
- `GET /{code}` route returning 302/404/403/410 without JSON envelope.

Stage 3 backend management is complete. The implementation includes:

- API contracts and handlers for administrator login, profile, and short-link management.
- RPC contracts and logic for administrator authentication and short-link management.
- JWT creation and validation in `link-api`.
- Bearer-token middleware for protected management routes.
- Stable API response envelopes for success and error responses.
- Domain validation for administrator status, URLs, short codes, expiration, pagination, and status values.
- 6-character base62 automatic short-code generation using `crypto/rand`.
- Custom short-code validation for 3-32 characters containing ASCII letters, digits, `_`, and `-`.
- Duplicate custom short-code conflict handling.
- Soft delete through `deleted_at`.
- Local HTTP smoke request assets under `tests/httpyac/`.

## Stage 3 Scope

Included:

- Administrator login.
- Authenticated administrator profile.
- JWT creation and validation for management APIs.
- Short-link creation.
- Short-link listing.
- Short-link detail retrieval.
- Short-link mutable field updates.
- Short-link soft delete.

Excluded:

- Redirect serving through `GET /{code}`.
- Redis redirect cache lookup, backfill, or invalidation.
- Visit recording.
- Statistics APIs.
- Management UI.
- Additional migrations beyond the current administrator and short-link foundation schema.

Public contract decisions:

- JWT uses local service configuration fields `Auth.Secret` and `Auth.TokenTTLSeconds`.
- Example configuration documents placeholder auth fields; local ignored configuration stores machine-local values.
- Auto-generated short codes are 6-character base62 values generated with `crypto/rand`.
- Custom short codes must be 3-32 characters and may contain ASCII letters, digits, `_`, and `-`.
- Reserved words are rejected.
- Duplicate custom short codes return `CONFLICT`.
- The short code is immutable after creation.
- Soft delete uses `deleted_at` and hides deleted links from normal management reads.

## Implemented HTTP API

The HTTP contracts are documented in `docs/06-api-design.md` and implemented in `services/link-api/`.

- `POST /admin/login`
- `GET /admin/profile`
- `POST /admin/links`
- `GET /admin/links`
- `GET /admin/links/{id}`
- `PATCH /admin/links/{id}`
- `DELETE /admin/links/{id}`
- `GET /admin/links/{id}/stats`

Implementation notes:

- `GET /healthz` and `GET /readyz` remain available.
- `GET /{code}` returns 302/404/403/410 and is wrapped by `AnalyticsMiddleware`.
- All management routes except login require Bearer token authentication.
- API logic owns HTTP parsing, response envelopes, JWT creation, and JWT validation.
- RPC logic owns credential verification, shared validation rules, and MySQL-backed business data.

## Implemented RPC API

The RPC contracts are documented in `docs/06-api-design.md` and implemented in `services/link-rpc/`.

- `AuthenticateAdmin`
- `GetAdminProfile`
- `CreateShortLink`
- `ListShortLinks`
- `GetShortLink`
- `UpdateShortLink`
- `DeleteShortLink`
- `ResolveShortLink`
- `RecordVisit`
- `GetLinkStats`

Implementation notes:

- The existing `Check` readiness method remains available.
- `ResolveShortLink` uses Redis cache lookup and status/expiry validation.
- `RecordVisit` hashes the raw IP server-side and inserts a `visit_event` row, then upserts `link_daily_stat`.
- `GetLinkStats` validates the date range and queries `link_daily_stat` between the given dates.
- Proto `go_package` uses an absolute import path without an explicit Go package alias.
- goctl-generated package names, client directories, and exported client identifiers are the source of truth.
- The services use generated go-zero `ServiceContext` wiring.
- No dependency injection framework is introduced.

## Generation Boundary

Stage 5 contracts have been generated. Future API or RPC contract changes must update the source contract
first, then run go-zero generation from the repository root. Do not handwrite generated go-zero service files.

API generation must use `--style gozero` to match the existing camelCase file naming convention:

```bash
goctl api go \
  --api services/link-api/api/link.api \
  --dir services/link-api \
  --style gozero
```

RPC generation command:

```bash
goctl rpc protoc services/link-rpc/proto/link/v1/link.proto \
  --go_out=services/link-rpc/pb \
  --go_opt=paths=source_relative \
  --go-grpc_out=services/link-rpc/pb \
  --go-grpc_opt=paths=source_relative \
  --zrpc_out=services/link-rpc \
  --proto_path=services/link-rpc/proto \
  -m
```

API generation should use the repository's existing go-zero API service layout and preserve current local
service configuration patterns.

## Local Smoke Flow

Run the following after entering the local development environment if needed:

```bash
make infra-up
make migrate-up
make run-rpc
make run-api
```

Create local httpyac variables:

```bash
cp tests/httpyac/http-client.example.env.json tests/httpyac/http-client.private.env.json
```

Then verify health and management flows:

```bash
httpyac send tests/httpyac/health.http -e local --all
httpyac send tests/httpyac/admin.http -e local -n login
```

After login, copy the returned token into `tests/httpyac/http-client.private.env.json`. After creating a
link, copy the returned ID into `linkId` before running detail, update, or delete requests.

## Current Verification

Run before merging changes:

```bash
make test
golangci-lint run ./...
docker compose -f deploy/docker-compose.infra.yml config --quiet
go test -race ./...
```

Run full-stack integration tests (requires `make infra-up && make migrate-up` first):

```bash
make test-e2e
```

## Stage 5 Acceptance

Stage 5 analytics implementation is accepted when:

- `GET /{code}` on an active link fires a `RecordVisit` RPC call asynchronously.
- `visit_event` table receives one row per redirect.
- `link_daily_stat` row for today shows `pv: 1` after the first redirect.
- `GET /admin/links/{id}/stats` returns the daily stats in the standard envelope.
- Analytics RPC failure does not block or error the redirect response.

## Stage 6 Acceptance

Stage 6 management UI is accepted when:

- Administrator can log in and is redirected to the links list.
- Unauthenticated access redirects to the login page.
- Administrator can create, list, view, edit, and soft-delete short links.
- Disabled links return 403 when accessed via short URL.
- Stats view displays daily PV/UV chart and summary cards for a given date range.
- Short link codes proxy correctly to link-api (`/:code` → 302) through the Vite dev server.

## Stage 7 Progress

Stage 7 testing, observability, and security hardening is complete.

Completed:

- Reserved short codes extended to include `metrics` and `static` in `link-rpc/internal/domain/validation.go`.
- `MaxConns: 1000` and `MaxBytes: 1048576` added to `link-api` `rest.RestConf` (configuration only).
- `IPRateLimitMiddleware` on `GET /{code}`: 20 req/s per IP using go-zero `limit.PeriodLimit` backed by Redis.
- `LoginRateLimitMiddleware` on `POST /admin/login`: 10 req/min per IP, reuses `IPRateLimitMiddleware`.
- `services/link-api/pkg/httputil.ExtractClientIP`: shared IP extraction utility (prefers `X-Forwarded-For`).
- Structured logging in `link-api`:
  - `loginlogic.go`: attempt, failed (with reason), success (with admin_id).
  - `authmiddleware.go`: authenticated (with admin_id and username).
  - `redirectlogic.go`: rpc failed (with code and error), resolved (with code and url).
- Structured logging in `link-rpc`:
  - `resolveshortlinklogic.go`: hit/miss/disabled/expired/error, each with code and result fields.
  - `authenticateadminlogic.go`: failed (username), success (admin_id, username).
  - `recordvisitlogic.go`: upsert daily stat failed with link_id, code, stat_date, error.
- Prometheus metrics in `link-rpc/internal/metrics/metrics.go`:
  - `zerolink_redirect_requests_total{result}` (hit/miss/disabled/expired/error).
  - `zerolink_analytics_events_total{result}` (success/error).
- Prometheus `/metrics` HTTP endpoints: link-api port 9100, link-rpc port 9101.
- Unit tests for `AnalyticsMiddleware`, `IPRateLimitMiddleware`, `LoginRateLimitMiddleware` using interface stubs and `go.uber.org/goleak`.
- Full-stack integration tests in `tests/integration/` covering health, auth, link CRUD, redirect (active/missing/disabled/expired), and analytics (RecordVisit + GetLinkStats polling). Services started as OS subprocesses by `TestMain`; run via `make test-e2e`.

## Stage 8 Progress

Stage 8 cache optimization and data lifecycle is complete.

Completed:

- Cuckoo filter package in `services/link-rpc/pkg/filter/` (thread-safe, wraps `github.com/seiflotfy/cuckoofilter`).
- On `NewServiceContext`: subscribe to `zl:code:created` Pub/Sub channel, batch-load all `short_link` codes into filter, start cleanup runner.
- `ResolveShortLink` fast path: if `CodeFilter != nil && !CodeFilter.Lookup(code)`, return `NotFound` immediately; `FilterRequestsTotal.Inc("miss")`.
- `CreateShortLink`: after successful insert, `CodeFilter.Insert(link.Code)` and `Redis.PublishCtx("zl:code:created", link.Code)`; publish failure is non-fatal.
- `CreateShortLink`: for custom codes, check `ReservedCodeModel.Exists` after `FindOneByCode` returns `ErrNotFound`; return `AlreadyExists` if reserved.
- Migration `000003_data_retention`: `short_link_archive` (non-auto-increment `id`) and `reserved_code` (VARCHAR primary key).
- `reserved_code` model generated by goctl `--style gozero`; `Exists` method added in `reservedcodemodel.go`.
- Cleanup runner: `visit_event` (90d), `short_link` soft-deleted (365d → archive → reserve), `link_daily_stat` (730d); all via batch deletes or per-row archival.
- New metrics: `zerolink_filter_requests_total{result}` and `zerolink_cleanup_deleted_rows_total{table}`.
- `Retention.*` and `Cuckoo.*` config sections with zero-value defaults in `NewServiceContext`.
- Unit tests: filter (basic + concurrent), runner cancellation, reserved-code conflict, filter miss/hit in resolve.

## Post-Stage-8 Quality Improvements

Completed after Stage 8 without advancing the feature stage boundary.

### Startup Config Validation

- New package `services/link-api/pkg/configvalidator` wraps `go-playground/validator/v10`.
- `MustValidate(cfg any)` validates struct tags and calls `logx.Severef` on failure — the process exits before serving traffic.
- Called in both `services/link-api/link.go` and `services/link-rpc/link.go` immediately after `conf.MustLoad`.
- `Auth.Secret`: `required,min=32`; `Auth.TokenTTLSeconds`: `required,gt=0`; `Analytics.IPSalt`: `required`.
- All config fields with `validate:` tags also carry an explicit `json:` tag to work around go-zero's `usingDifferentKeys` field-skip behaviour.
- Removed runtime `validateConfig()` path from `services/link-api/internal/auth/token.go`; the startup check is the single validation point.

### Error Handler Improvements

- go-zero `httpx.Parse` failures (missing required field, type mismatch, malformed JSON body) now return `BAD_REQUEST 400` instead of `INTERNAL 500`. Handled in `apierror.isParseError` by inspecting the error message strings produced by go-zero.
- gRPC `codes.Unavailable` now maps to `SERVICE_UNAVAILABLE 503` instead of `INTERNAL 500`.
- See `services/link-api/internal/apierror/error.go` for the full mapping table.

## Stage 9 Progress

Stage 9 full Docker Compose deployment is complete.

Completed:

- `services/link-api/Dockerfile` — multi-stage build (`golang:1.26-alpine` builder, `alpine:3.21` runtime), `CGO_ENABLED=0`, non-root `appuser`, exposes 8080 and 9100.
- `services/link-rpc/Dockerfile` — same pattern, exposes 9090 and 9101.
- `etc/link-api.compose.yaml` — `Host: 0.0.0.0`, `Prometheus.Host: 0.0.0.0`, `Mode: console`, `Encoding: json`, all values via `${VAR}` env substitution.
- `etc/link-rpc.compose.yaml` — `ListenOn: 0.0.0.0:9090`, `loc=UTC` in MySQL DSN, same env var pattern.
- `.env.compose.example` — documents all required and optional env vars; secrets (`LINK_API_AUTH_SECRET`, `LINK_RPC_IP_SALT`) have no defaults. Copy to `.env.compose.local` (gitignored) for local use.
- `deploy/docker-compose.yml` — five services: mysql, redis, migrate (one-shot `ghcr.io/golang-migrate/migrate`), link-rpc, link-api. Only port 8080 published to host. `service_completed_successfully` condition on migrate ensures schema is applied before link-rpc starts.
- `Makefile` — added `compose-build`, `compose-up`, `compose-down`, `compose-logs` targets.
- `docs/11-deployment-design.md` — Stage 9 Implementation section added.

## Next Handoff

Stage 9 is the final planned implementation stage. Future work is tracked in `docs/15-future-roadmap.md`.

## Git And Commit Rules

- Do not make task changes directly on `main`.
- Before creating any commit, show the intended staged changes and commit message, then wait for explicit user confirmation.
- Use Conventional Commits without emoji.
