# Implementation Plan

## Goal

Track the current zero-link implementation state and define the next safe handoff point. The project has moved past the Stage 1 repository foundation and completed the Stage 2 API/RPC skeleton. It is now preparing the Stage 3 administrator and short-link management implementation.

This document is currently a Stage 3 specification handoff. It intentionally does not add short-link business routes, business RPC methods, additional migrations, redirect behavior, admin UI, analytics, or handwritten generated code.

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

Stage 2 is implemented. The current transport skeleton includes:

- go-zero API service skeleton under `services/link-api/`.
- go-zero RPC service skeleton under `services/link-rpc/`.
- API contract at `services/link-api/api/link.api`.
- RPC protobuf contract at `services/link-rpc/proto/link/v1/link.proto`.
- Generated protobuf and gRPC code under `services/link-rpc/pb/link/v1/`.
- RPC client package under the goctl-generated `services/link-rpc/linkservice/` directory.
- API route contracts for `GET /healthz` and `GET /readyz` only.
- RPC readiness method `LinkService.Check` only.
- API readiness logic that calls the RPC readiness method.
- RPC readiness logic that validates configured MySQL and Redis endpoints and performs simple TCP connectivity checks.
- Unit coverage for RPC readiness validation of missing dependency endpoints.

Stage 3 foundation work has started. The current foundation includes:

- `golang-migrate` commands for local schema application and rollback.
- Initial migration files for `admin_user` and `short_link`.
- A seeded local/dev administrator account for migration smoke testing.
- Generated go-zero MySQL models for `admin_user` and `short_link`.
- RPC service-context wiring for the generated models.

## Stage 3 Specification Decisions

Stage 3 covers administrator authentication and short-link management only.

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
- Additional migrations.

Public contract decisions:

- JWT uses local service configuration fields `Auth.Secret` and `Auth.TokenTTLSeconds`.
- Example configuration documents placeholder auth fields; local ignored configuration stores machine-local values.
- Auto-generated short codes are 6-character base62 values generated with `crypto/rand`.
- Custom short codes must be 3-32 characters and may contain ASCII letters, digits, `_`, and `-`.
- Reserved words are rejected.
- Duplicate custom short codes return `CONFLICT`.
- The short code is immutable after creation.
- Soft delete uses `deleted_at` and hides deleted links from normal management reads.

## Stage 3 HTTP Handoff

The next implementation pass should update `services/link-api/api/link.api` to add only the Stage 3 management HTTP contracts documented in `docs/06-api-design.md`:

- `POST /admin/login`
- `GET /admin/profile`
- `POST /admin/links`
- `GET /admin/links`
- `GET /admin/links/{id}`
- `PATCH /admin/links/{id}`
- `DELETE /admin/links/{id}`

Implementation notes:

- Keep `GET /healthz` and `GET /readyz`.
- Do not add `GET /{code}`.
- Do not add `GET /admin/links/{id}/stats`.
- Management responses use the stable response envelope from `docs/06-api-design.md`.
- All management routes except login require Bearer token authentication.
- API logic creates and validates JWTs; RPC logic authenticates credentials and manages MySQL-backed business data.

## Stage 3 RPC Handoff

The next implementation pass should update `services/link-rpc/proto/link/v1/link.proto` to add only these Stage 3 RPC methods:

- `AuthenticateAdmin`
- `GetAdminProfile`
- `CreateShortLink`
- `ListShortLinks`
- `GetShortLink`
- `UpdateShortLink`
- `DeleteShortLink`

Implementation notes:

- Keep the existing `Check` readiness method.
- Do not add `ResolveShortLink`, `RecordVisit`, or `GetLinkStats`.
- Keep proto `go_package` absolute without an explicit Go package alias.
- Accept goctl-generated package names, client directories, and exported client identifiers as the source of truth.
- Use generated go-zero `ServiceContext` wiring.
- Do not introduce a dependency injection framework.

## Generation Handoff

After this specification is reviewed, the next implementation pass must explicitly run go-zero generation from the repository root. Do not handwrite generated go-zero service files.

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

API generation should use the repository's existing go-zero API service layout and preserve current local service configuration patterns.

## Stage 3 Implementation Sequence

Recommended next implementation order:

1. Create a task branch from `main`.
2. Update API and proto contracts for Stage 3 management only.
3. Run approved goctl generation from the repository root.
4. Add auth configuration fields to the API config and example config.
5. Add JWT creation and validation helpers with unit tests.
6. Add password verification and administrator authentication logic with unit tests.
7. Add short-code generation and validation helpers with unit tests.
8. Add URL, status, expiration, pagination, and reserved-word validation with unit tests.
9. Implement management RPC logic against the existing generated MySQL models.
10. Implement API handlers and middleware using generated service context wiring.
11. Add integration tests guarded by the `integration` build tag.
12. Run `make test`.
13. Validate the local infrastructure Compose config.
14. Run Stage 3 smoke tests against local infrastructure.

## Current Verification

Run the following after entering the local development environment if needed:

```bash
make test
docker compose -f deploy/docker-compose.infra.yml config --quiet
```

For a local service smoke test:

```bash
make infra-up
make migrate-up
make run-rpc
make run-api
```

Then verify:

```bash
curl http://127.0.0.1:8080/healthz
curl http://127.0.0.1:8080/readyz
```

Expected current behavior:

- `/healthz` reports the API process is healthy.
- `/readyz` reports the API and RPC are ready when MySQL and Redis are reachable.
- Stage 3 management endpoints are not available until the next implementation pass explicitly updates contracts and runs goctl generation.

## Stage 3 Acceptance

Stage 3 management implementation is accepted when:

- `POST /admin/login` succeeds for an active administrator.
- Login rejects invalid credentials and inactive administrators.
- Authenticated profile retrieval works.
- Management APIs reject missing or invalid Bearer tokens.
- Short-link creation supports generated and valid custom codes.
- Invalid URLs, reserved codes, invalid codes, invalid expiration values, and duplicate custom codes are rejected with stable envelope errors.
- Short-link listing supports pagination and management filters.
- Short-link detail retrieval hides missing and soft-deleted links.
- Short-link update changes mutable fields and keeps code immutable.
- Short-link delete soft deletes and hides the link from normal management reads.
- Unit tests and integration tests cover the Stage 3 scenarios in `docs/12-testing-strategy.md`.
- Redirect serving, Redis redirect cache, analytics, and management UI remain deferred.

## Git And Commit Rules

- Do not make task changes directly on `main`.
- Before creating any commit, show the intended staged changes and commit message, then wait for explicit user confirmation.
- Use Conventional Commits without emoji.
