# Local Development

## Development Mode

Local development runs infrastructure in Docker Compose and Go services on the host machine.

This gives fast rebuilds, easy debugger attachment, and predictable MySQL/Redis dependencies.

## Planned Files

- `deploy/docker-compose.infra.yml`: starts MySQL and Redis.
- `.env.example`: documents local environment values.
- `.env.local`: machine-local environment values ignored by Git.
- `etc/link-api.example.yaml`: committed API service config example.
- `etc/link-api-local.yaml`: machine-local API service config ignored by Git.
- `etc/link-rpc.example.yaml`: committed RPC service config example.
- `etc/link-rpc-local.yaml`: machine-local RPC service config ignored by Git.
- `Makefile`: wraps common development commands.

## Dependency Services

`deploy/docker-compose.infra.yml` should expose:

- MySQL on `127.0.0.1:3306`.
- Redis on `127.0.0.1:6379`.
- Optional Adminer or equivalent on a local-only port.

Use named volumes so database state survives container restarts.

## Common Commands

Planned commands:

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

The exact commands will be implemented after the service structure exists.

## Local Verification Flow

1. Start MySQL and Redis.
2. Run migrations.
3. Start `link-rpc` locally.
4. Start `link-api` locally.
5. Log in as administrator.
6. Create a short link.
7. Open `http://localhost:{api_port}/{code}`.
8. Confirm redirect status and destination.
9. View link statistics.

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
