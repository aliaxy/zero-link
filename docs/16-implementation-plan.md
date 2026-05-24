# Stage 1 Skeleton Completion Plan

## Goal

Define and verify the completed skeleton-only engineering foundation for zero-link without generating go-zero service code, defining API/RPC contracts, creating business schema, or implementing short-link behavior. Stage 1 keeps the repository documented, initialized, and ready for later goctl-based work.

## Source Documents

This plan is derived from:

- `docs/00-project-overview.md`
- `docs/03-system-architecture.md`
- `docs/10-local-development.md`
- `docs/14-milestones.md`

## Decisions

- Module path: `github.com/aliaxy/zero-link`.
- Local dependency mode: Docker Compose runs MySQL and Redis.
- Local service mode: Go services run on the host machine.
- Service shape: reserve `link-api` and `link-rpc` directories only.
- Stage 1 transport: no API routes, no protobuf methods, and no generated HTTP/RPC code.
- Service discovery: direct local target, no required etcd.
- Admin UI: reserve `web/admin` for a future embedded React UI.
- Dependency injection: use go-zero generated service contexts later, no DI framework.
- Stage 2 owns generated go-zero service skeletons, health checks, API contracts, protobuf methods, and the first API-to-RPC local call.

## Completed Task 1: Repository Baseline

The project baseline includes:

- `.gitignore`
- `.env.example`
- `Makefile`
- `README.md`
- `go.mod` using `go mod init github.com/aliaxy/zero-link`

Acceptance:

- The repository has standard ignored files.
- Local configuration values are documented without real secrets.
- Common commands are available through `make`.

## Completed Task 2: Service Layout

The reserved layout includes:

- `services/link-api/api`
- `services/link-rpc/proto`
- `etc`
- `migrations`
- `web/admin`

Acceptance:

- No handwritten API/RPC service skeleton exists.
- No API/RPC contracts exist yet.
- Future generated code will own the service entrypoints and service contexts.

## Completed Task 3: Local Infrastructure

`deploy/docker-compose.infra.yml` provides:

- MySQL exposed on `127.0.0.1:3306`.
- Redis exposed on `127.0.0.1:6379`.
- Named volumes for local data.
- Health checks for both services.

Acceptance:

- `make infra-up` starts MySQL and Redis.
- `make infra-down` stops them.

## Completed Task 4: Configuration

Committed examples:

- `etc/link-api.example.yaml`
- `etc/link-rpc.example.yaml`

Ignored local copies for machine-specific values:

- `.env.local`
- `etc/link-api-local.yaml`
- `etc/link-rpc-local.yaml`

Local copies must be ignored by Git.

The first implementation can use simple environment variables with these YAML files as documented configuration targets. Full YAML parsing can be added when go-zero generated config structs are introduced.

Acceptance:

- Example configuration files are committed.
- Local configuration files exist for development but are not tracked.
- API defaults to `127.0.0.1:8080`.
- RPC defaults to `127.0.0.1:9090`.
- RPC dependency defaults point to local MySQL and Redis.

## Completed Task 5: Migration Placeholder

Only a migration placeholder exists:

- `migrations/README.md`

Acceptance:

- No business schema exists in Stage 1.
- The future migration location is clear.

## Completed Task 6: Future Service Contract Placeholder

Keep placeholders ready for later generated service behavior:

- `services/link-api/api/README.md`
- `services/link-rpc/proto/README.md`

Acceptance:

- No handwritten implementation is added before goctl generation.
- No API route or RPC method is committed in Stage 1.
- Later generated handlers and logic will implement service behavior.

## Deferred To Stage 2

The following items are intentionally deferred until the user explicitly starts the implementation stage:

- go-zero API/RPC generation.
- API route contracts.
- Protobuf RPC methods.
- Health and readiness endpoints.
- API-to-RPC local call verification.
- Business migrations and model code.
- Short-link business logic.

## Current Verification

Run:

```bash
make test
docker compose -f deploy/docker-compose.infra.yml config --quiet
```

If no Go packages exist yet, `make test` should report that the skeleton foundation is ready.

Acceptance:

- No handwritten service code exists.
- No API/RPC contracts or business migrations exist yet.
- Current verification commands pass.
- Documentation contains no placeholder text.
