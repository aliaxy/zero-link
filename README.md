# zero-link

zero-link is a go-zero-oriented short-link system. Stages 1–5 are complete.

## Current Stage

Stage 5 adds analytics on top of the Stage 4 redirect and Stage 3 management foundation:

- MySQL and Redis are provided by `deploy/docker-compose.infra.yml`.
- `link-api` exposes `GET /healthz`, `GET /readyz`, `GET /{code}` redirect, and management + analytics APIs.
- `link-rpc` exposes readiness, administrator auth, short-link management, redirect resolution, visit recording, and stats query.
- Stage 3 migrations create `admin_user` and `short_link`.
- Stage 5 migrations create `visit_event` and `link_daily_stat`.
- `GET /{code}` returns 302/404/403/410; `AnalyticsMiddleware` fires a non-blocking `RecordVisit` RPC goroutine on every 302.
- `GET /admin/links/{id}/stats` returns daily PV/UV for a short link within a date range.
- Administrator login and profile APIs use JWT-based authentication.
- Authenticated short-link create, list, detail, update, and soft-delete management APIs are implemented.
- Local HTTP smoke request assets are available under `tests/httpyac/`.
- Management UI is deferred to Stage 6.
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

Apply the current Stage 5 schema when local MySQL is running:

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
