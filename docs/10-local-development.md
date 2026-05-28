# Local Development

## Development Mode

Local development runs infrastructure in Docker Compose and Go services on the host machine.

This gives fast rebuilds, easy debugger attachment, and predictable MySQL/Redis dependencies.

## Local Files

- `deploy/docker-compose.infra.yml`: starts MySQL and Redis.
- `.env.example`: documents local environment values.
- `.env.local`: machine-local environment values ignored by Git.
- `etc/link-api.example.yaml`: committed API service config example.
- `etc/link-api-local.yaml`: machine-local API service config ignored by Git.
- `etc/link-rpc.example.yaml`: committed RPC service config example.
- `etc/link-rpc-local.yaml`: machine-local RPC service config ignored by Git.
- `tests/httpyac/*.http`: local HTTP smoke requests.
- `tests/httpyac/http-client.example.env.json`: committed httpyac variable example.
- `tests/httpyac/http-client.private.env.json`: machine-local httpyac variables ignored by Git.
- `Makefile`: wraps common development commands.

## Dependency Services

`deploy/docker-compose.infra.yml` should expose:

- MySQL on `127.0.0.1:3306`.
- Redis on `127.0.0.1:6379`.
- Optional Adminer or equivalent on a local-only port.

Use named volumes so database state survives container restarts.

## Common Commands

```bash
make infra-up
make infra-down
make migrate-up
make migrate-down
make run-rpc
make run-api
make test
make test-integration
```

`make run-rpc` starts the RPC service with `etc/link-rpc-local.yaml`.
`make run-api` starts the API service with `etc/link-api-local.yaml`.

## Migration Workflow

Stage 3 foundation uses a `golang-migrate` workflow.

Intended local behavior:

- `make migrate-up` applies all pending migrations to the local MySQL database.
- `make migrate-down` rolls back the latest local migration step.
- Migration files live under `migrations/` and use versioned up/down SQL.
- The first Stage 3 migration creates `admin_user` and `short_link`.
- A local/dev seed administrator may be inserted by migration SQL so login can be tested without manual database edits.

Do not run migrations before `make infra-up` has started MySQL.

## Local Verification Flow

1. Start MySQL and Redis with `make infra-up`.
2. Run `make migrate-up` to apply the Stage 3 foundation schema.
3. Start `link-rpc` locally with `make run-rpc`.
4. Start `link-api` locally with `make run-api`.
5. Open `http://127.0.0.1:8080/healthz`.
6. Open `http://127.0.0.1:8080/readyz`.
7. Confirm `healthz` succeeds when the API process is alive.
8. Confirm `readyz` succeeds only when API, RPC, MySQL, and Redis are reachable.
9. Log in with the seeded local administrator and create a short link through the management API.

## HTTP Smoke Requests

The repository keeps local smoke requests under `tests/httpyac/`.

Create a private httpyac environment file before running authenticated requests:

```bash
cp tests/httpyac/http-client.example.env.json tests/httpyac/http-client.private.env.json
```

Use the `local` environment when running requests:

```bash
httpyac send tests/httpyac/health.http -e local --all
httpyac send tests/httpyac/admin.http -e local -n login
```

After login, copy the returned bearer token into `tests/httpyac/http-client.private.env.json`.
After creating a link, copy the returned link ID into `linkId` before running detail, update, or delete
requests.

## Configuration Rules

- Committed configuration files are examples only.
- Local configuration files are copied from examples and ignored by Git.
- Local service config uses `127.0.0.1` for MySQL and Redis.
- Compose deployment config uses service names.
- Secrets are read from environment variables or local ignored files.
- `.env.example` documents required values but must not contain real secrets.

Create local configuration before running services:

```bash
cp .env.example .env.local
cp etc/link-api.example.yaml etc/link-api-local.yaml
cp etc/link-rpc.example.yaml etc/link-rpc-local.yaml
```

Only edit the local files for machine-specific values.

## Tooling Note

The repository includes a Nix development shell with Go, goctl, protobuf tooling, and lint tools. If `go` or `goctl` are not available in a normal shell, enter the Nix development environment before running local commands.
