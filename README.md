# zero-link

zero-link is a go-zero-oriented short-link system. Stage 2 API/RPC skeleton development is complete, and the project is now aligning the Stage 3 foundation: initial administrator and short-link schema, generated MySQL models, and local migration tooling.

## Current Stage

Stage 3 foundation keeps the Stage 2 transport skeleton in place while adding database groundwork for upcoming administrator and short-link management:

- MySQL and Redis are provided by `deploy/docker-compose.infra.yml`.
- `link-api` exposes `GET /healthz` and `GET /readyz`.
- `link-rpc` exposes a minimal readiness RPC used by the API.
- Stage 3 migrations create `admin_user` and `short_link`.
- Generated MySQL models are wired into `link-rpc` service context for future use.
- Admin APIs, business RPC methods, redirect APIs, cache behavior, analytics, UI, and short-link business logic are still deferred.
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

Apply the current Stage 3 foundation schema when local MySQL is running:

```bash
make migrate-up
```

Short-link business implementation is reserved for the next Stage 3 implementation pass.

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
