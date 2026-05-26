# Implementation Plan

## Goal

Track the current zero-link implementation state and define the next safe handoff point. The project has moved past the Stage 1 repository foundation and completed the Stage 2 API/RPC skeleton. It is now in Stage 3 foundation alignment.

This plan intentionally preserves the current Stage 3 foundation boundary: no short-link business APIs, business RPC methods, additional database migrations, redirect behavior, admin UI, or analytics should be added until the next implementation pass is explicitly started.

## Source Documents

This plan is derived from:

- `docs/00-project-overview.md`
- `docs/03-system-architecture.md`
- `docs/10-local-development.md`
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

## Stage 2 Decisions

- Keep service discovery local and direct for now.
- Keep Docker Compose limited to local dependencies only.
- Keep application services running on the host machine during local development.
- Keep readiness checks simple and schema-free.
- Use generated go-zero `ServiceContext` wiring.
- Keep proto `go_package` absolute without an explicit Go package alias.
- Accept goctl-generated package names, client directories, and exported client identifiers as the source of truth.
- Do not introduce a dependency injection framework.
- Do not add short-link business behavior in Stage 2.

## Stage 2 Acceptance

Stage 2 is accepted:

- `GET /healthz` returns a successful API health response.
- `GET /readyz` calls `link-rpc` and returns ready when RPC dependencies are reachable.
- `link-rpc` readiness reports unavailable dependencies without requiring any business schema.
- `make test` passes in a shell where Go is available.
- The local infrastructure Compose file validates.
- No admin, redirect, analytics, or short-link management API has been introduced.

## Current Verification

Run the following after entering the local development environment if needed:

```bash
make test
docker compose -f deploy/docker-compose.infra.yml config --quiet
```

For a local service smoke test:

```bash
make infra-up
make run-rpc
make run-api
```

Then verify:

```bash
curl http://127.0.0.1:8080/healthz
curl http://127.0.0.1:8080/readyz
```

Expected behavior:

- `/healthz` reports the API process is healthy.
- `/readyz` reports the API and RPC are ready when MySQL and Redis are reachable.

## Stage 3 Foundation Status

The active implementation stage is Stage 3: Administrator And Link Management. The first foundation slice has begun with migration tooling, initial schema, generated models, and model wiring.

Before writing Stage 3 business behavior, create a focused Stage 3 implementation plan that covers:

- Business database schema for administrators and links.
- Migration workflow and local migration verification.
- Administrator authentication scope.
- Link creation, listing, detail, update, and soft delete APIs.
- RPC methods required by those API operations.
- Cache invalidation boundaries that will matter in Stage 4.
- Unit and integration test strategy.

Stage 3 should not include redirect serving, Redis redirect caching, analytics, or the management UI. Those remain reserved for later milestones.

## Stage 3 Next Task Recommendation

Continue Stage 3 with the administrator and short-link management implementation slice:

1. Keep the existing `docs/04-storage-design.md`, `docs/06-api-design.md`, and `docs/12-testing-strategy.md` as the source of truth.
2. Add administrator authentication, JWT handling, management API contracts, business RPC methods, and tests in a focused implementation pass.
3. Do not add redirect serving, Redis redirect caching, analytics, or management UI during Stage 3 management implementation.

This keeps the business implementation aligned with the existing milestone boundaries and avoids pulling redirect/cache behavior forward too early.
