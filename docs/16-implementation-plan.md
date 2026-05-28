# Implementation Plan

## Goal

Track the current zero-link implementation state and define the next safe handoff point. Stages 1–4 are
complete. The project is ready to begin Stage 5 analytics.

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
- RPC readiness method `LinkService.Check`.
- API readiness logic that calls the RPC readiness method.
- RPC readiness logic that validates configured MySQL and Redis endpoints with simple connectivity checks.

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

The Stage 3 management HTTP contracts are documented in `docs/06-api-design.md` and implemented in
`services/link-api/`.

- `POST /admin/login`
- `GET /admin/profile`
- `POST /admin/links`
- `GET /admin/links`
- `GET /admin/links/{id}`
- `PATCH /admin/links/{id}`
- `DELETE /admin/links/{id}`

Implementation notes:

- `GET /healthz` and `GET /readyz` remain available.
- `GET /{code}` is implemented and returns 302/404/403/410.
- `GET /admin/links/{id}/stats` is still deferred.
- All management routes except login require Bearer token authentication.
- API logic owns HTTP parsing, response envelopes, JWT creation, and JWT validation.
- RPC logic owns credential verification, shared validation rules, and MySQL-backed business data.

## Implemented RPC API

The Stage 3 RPC contracts are documented in `docs/06-api-design.md` and implemented in `services/link-rpc/`.

- `AuthenticateAdmin`
- `GetAdminProfile`
- `CreateShortLink`
- `ListShortLinks`
- `GetShortLink`
- `UpdateShortLink`
- `DeleteShortLink`

Implementation notes:

- The existing `Check` readiness method remains available.
- `ResolveShortLink` is implemented with Redis cache lookup and status/expiry validation.
- `RecordVisit` and `GetLinkStats` remain deferred.
- Proto `go_package` uses an absolute import path without an explicit Go package alias.
- goctl-generated package names, client directories, and exported client identifiers are the source of truth.
- The services use generated go-zero `ServiceContext` wiring.
- No dependency injection framework is introduced.

## Generation Boundary

Stage 3 contracts have been generated. Future API or RPC contract changes must update the source contract
first, then run go-zero generation from the repository root. Do not handwrite generated go-zero service files.

RPC generation command:

```bash
goctl rpc protoc services/link-rpc/proto/link/v1/link.proto \
  --go_out=services/link-rpc/pb \
  --go_opt=paths=source_relative \
  --go-grpc_out=services/link-rpc/pb \
  --go-grpc_opt=paths=source_relative \
  --zrpc_out=services/link-rpc \
  --proto_path=services/link-rpc/proto
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

## Stage 4 Acceptance

Stage 4 redirect and cache implementation is accepted when:

- Active links redirect with `302 Found`.
- Missing and soft-deleted links return `404 Not Found`.
- Disabled links return `403 Forbidden`.
- Expired links return `410 Gone`.
- Redis cache is populated on first resolve and invalidated on update and soft delete.
- Analytics and management UI remain deferred.

## Next Handoff

Stage 5 analytics work:

- Add `RecordVisit` RPC method and visit event storage.
- Add `GetLinkStats` RPC method for daily stat aggregation.
- Add `GET /admin/links/{id}/stats` HTTP endpoint.
- Keep management UI deferred to Stage 6.

## Git And Commit Rules

- Do not make task changes directly on `main`.
- Before creating any commit, show the intended staged changes and commit message, then wait for explicit user confirmation.
- Use Conventional Commits without emoji.
