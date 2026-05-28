# zero-link

zero-link is a go-zero-oriented short-link system. Stage 2 API/RPC skeleton development is complete, and the Stage 3 backend management implementation is now in place.

## Current Stage

Stage 3 keeps the Stage 2 transport skeleton in place while adding administrator authentication and MySQL-backed short-link management:

- MySQL and Redis are provided by `deploy/docker-compose.infra.yml`.
- `link-api` exposes `GET /healthz` and `GET /readyz`.
- `link-rpc` exposes readiness checks used by the API.
- Stage 3 migrations create `admin_user` and `short_link`.
- Generated MySQL models are wired into `link-rpc` service context.
- Administrator login and profile APIs are implemented with JWT-based authentication.
- Authenticated short-link create, list, detail, update, and soft-delete management APIs are implemented.
- Local HTTP smoke request assets are available under `tests/httpyac/`.
- Redirect APIs, Redis redirect cache behavior, analytics, and the management UI are still deferred.
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

Apply the current Stage 3 schema when local MySQL is running:

```bash
make migrate-up
```

Run local HTTP smoke requests from `tests/httpyac/` after the API, RPC, and local dependencies are running.

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
