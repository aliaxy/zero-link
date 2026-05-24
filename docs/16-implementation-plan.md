# Stage 1 Skeleton Plan

## Goal

Build the skeleton-only engineering foundation for zero-link without generating go-zero service code, defining API/RPC contracts, creating business schema, or implementing short-link behavior. Stage 1 must make the repository documented, initialized, and ready for later goctl-based work.

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

## Task 1: Repository Baseline

Create the project baseline files:

- `.gitignore`
- `.env.example`
- `Makefile`
- `README.md`
- `go.mod` using `go mod init github.com/aliaxy/zero-link`

Acceptance:

- The repository has standard ignored files.
- Local configuration values are documented without real secrets.
- Common commands are available through `make`.

## Task 2: Service Layout

Create the initial layout:

- `services/link-api/api`
- `services/link-rpc/proto`
- `etc`
- `migrations`
- `web/admin`

Acceptance:

- No handwritten API/RPC service skeleton exists.
- No API/RPC contracts exist yet.
- Future generated code will own the service entrypoints and service contexts.

## Task 3: Local Infrastructure

Create `deploy/docker-compose.infra.yml` with:

- MySQL exposed on `127.0.0.1:3306`.
- Redis exposed on `127.0.0.1:6379`.
- Named volumes for local data.
- Health checks for both services.

Acceptance:

- `make infra-up` starts MySQL and Redis.
- `make infra-down` stops them.

## Task 4: Configuration

Create:

- `etc/link-api-local.yaml`
- `etc/link-rpc-local.yaml`

The first implementation can use simple environment variables with these YAML files as documented configuration targets. Full YAML parsing can be added when go-zero generated config structs are introduced.

Acceptance:

- API defaults to `127.0.0.1:8080`.
- RPC defaults to `127.0.0.1:9090`.
- RPC dependency defaults point to local MySQL and Redis.

## Task 5: Migration Placeholder

Create only a migration placeholder:

- `migrations/README.md`

Acceptance:

- No business schema exists in Stage 1.
- The future migration location is clear.

## Task 6: Future Service Contract Placeholder

Keep placeholders ready for later generated service behavior:

- `services/link-api/api/README.md`
- `services/link-rpc/proto/README.md`

Acceptance:

- No handwritten implementation is added before goctl generation.
- No API route or RPC method is committed in Stage 1.
- Later generated handlers and logic will implement service behavior.

## Task 7: Tests

Before generated code exists:

- `make test` should succeed without requiring Go packages.
- Contract syntax validation happens later when contracts exist.

After generated code exists, add unit tests for:

- generated handlers.
- service logic.
- dependency behavior.

## Task 8: Verification

Run:

```bash
gofmt -w .
go test ./...
```

If no Go packages exist yet, `make test` should report that the skeleton foundation is ready.

Acceptance:

- No handwritten service code exists.
- No API/RPC contracts or business migrations exist yet.
- Current verification commands pass.
- Documentation contains no placeholder text.
