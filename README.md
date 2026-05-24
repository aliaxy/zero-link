# zero-link

zero-link is a go-zero-oriented short-link system. The project is currently in Stage 2 API/RPC skeleton development: go-zero service skeletons exist for health/readiness, while short-link business behavior and business schema are intentionally deferred.

## Current Stage

Stage 2 builds the minimal API/RPC transport skeleton without short-link business implementation:

- MySQL and Redis are provided by `deploy/docker-compose.infra.yml`.
- `link-api` exposes `GET /healthz` and `GET /readyz`.
- `link-rpc` exposes a minimal readiness RPC used by the API.
- Business migrations, admin APIs, redirect APIs, and short-link logic are still deferred.
- `go.mod` is initialized with `github.com/aliaxy/zero-link`.

## Requirements

- Go 1.26 or newer.
- Docker with Compose support.
- Optional: Nix development shell from `flake.nix`.

If `go` is not available in your normal shell, enter the Nix development environment first.

## Local Development

Create local configuration from the committed examples:

```bash
cp .env.example .env.local
cp etc/link-api.example.yaml etc/link-api-local.yaml
cp etc/link-rpc.example.yaml etc/link-rpc-local.yaml
```

The `*.example.*` files are committed templates. The `*.local.*` files are for machine-local values and are ignored by Git.

```bash
make infra-up
```

Run the generated services locally:

```bash
make run-rpc
make run-api
```

Business migrations and short-link implementation are reserved for later stages.

Validate the current repository:

```bash
make test
```

## Documentation

Start with:

- `docs/00-project-overview.md`
- `docs/03-system-architecture.md`
- `docs/10-local-development.md`
- `docs/16-implementation-plan.md`
