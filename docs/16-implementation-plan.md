# Implementation Plan

## Goal

Track the current zero-link implementation state and define the next safe handoff point. The project has moved past the Stage 1 repository foundation and is now at the end of Stage 2 API/RPC skeleton work.

This plan intentionally preserves the current stage boundary: no short-link business APIs, business RPC methods, database migrations, redirect behavior, admin UI, or analytics should be added until the next implementation stage is explicitly started.

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
- Reserved directories for services, migrations, and future admin UI.
- Make targets for local infrastructure, service runs, tests, formatting, and linting.

Stage 2 is implemented. The current transport skeleton includes:

- go-zero API service skeleton under `services/link-api/`.
- go-zero RPC service skeleton under `services/link-rpc/`.
- API contract at `services/link-api/api/link.api`.
- RPC protobuf contract at `services/link-rpc/proto/link/v1/link.proto`.
- Generated protobuf and gRPC code under `services/link-rpc/pb/link/v1/`.
- RPC client package under `services/link-rpc/linkclient/`.
- API route contracts for `GET /healthz` and `GET /readyz` only.
- RPC readiness method `LinkService.Check` only.
- API readiness logic that calls the RPC readiness method.
- RPC readiness logic that validates configured MySQL and Redis endpoints and performs simple TCP connectivity checks.
- Unit coverage for RPC readiness validation of missing dependency endpoints.

## Stage 2 Decisions

- Keep service discovery local and direct for now.
- Keep Docker Compose limited to local dependencies only.
- Keep application services running on the host machine during local development.
- Keep readiness checks simple and schema-free.
- Use generated go-zero `ServiceContext` wiring.
- Keep proto `go_package` absolute with explicit alias `linkv1`.
- Normalize generated code imports to clear package aliases when needed.
- Do not introduce a dependency injection framework.
- Do not add short-link business behavior in Stage 2.

## Stage 2 Acceptance

Stage 2 is accepted when:

- `GET /healthz` returns a successful API health response.
- `GET /readyz` calls `link-rpc` and returns ready when RPC dependencies are reachable.
- `link-rpc` readiness reports unavailable dependencies without requiring any business schema.
- `make test` passes in a shell where Go is available.
- The local infrastructure Compose file validates.
- No business migrations exist.
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

## Stage 3 Entry Point

The next implementation stage is Stage 3: Administrator And Link Management.

Before writing Stage 3 code, create a focused Stage 3 plan that covers:

- Business database schema for administrators and links.
- Migration workflow and local migration verification.
- Administrator authentication scope.
- Link creation, listing, detail, update, and soft delete APIs.
- RPC methods required by those API operations.
- Cache invalidation boundaries that will matter in Stage 4.
- Unit and integration test strategy.

Stage 3 should not include redirect serving, Redis redirect caching, analytics, or the management UI. Those remain reserved for later milestones.

## Stage 3 First Task Recommendation

Start Stage 3 with a documentation-first design update for the persistence and management API slice:

1. Update `docs/04-storage-design.md` with the concrete administrator and link tables for the first business schema.
2. Update `docs/06-api-design.md` with the exact management API contracts.
3. Update `docs/12-testing-strategy.md` with the Stage 3 unit and integration coverage.
4. Only after the design is accepted, add migrations, RPC methods, API routes, generated service code, and tests.

This keeps the business implementation aligned with the existing milestone boundaries and avoids pulling redirect/cache behavior forward too early.
