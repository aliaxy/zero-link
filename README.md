# zero-link

zero-link is a go-zero-oriented short-link system. Stages 1–6 are complete. Stage 7 is in progress.

## Current Stage

Stage 7 adds observability, security hardening, and testing on top of the Stage 6 admin SPA foundation:

- IP rate limiting on `GET /{code}` (20 req/s per IP) and `POST /admin/login` (10 req/min per IP) via go-zero `PeriodLimit`.
- Structured logging across `link-api` (login, auth, redirect) and `link-rpc` (resolve, authenticate, analytics) using `logx.Infow` / `logx.Errorw`.
- Prometheus `/metrics` endpoints: `link-api` on port 9100, `link-rpc` on port 9101.
- Business metric counters `zerolink_redirect_requests_total{result}` and `zerolink_analytics_events_total{result}`.
- Unit tests for `AnalyticsMiddleware`, `IPRateLimitMiddleware`, and `LoginRateLimitMiddleware` with goroutine leak detection via `goleak`.
- `MaxConns: 1000` and `MaxBytes: 1048576` connection protection on `link-api`.
- Reserved short codes extended to include `metrics` and `static`.

Stage 6 added a Vue 3 + Vite + Element Plus admin SPA:

- MySQL and Redis are provided by `deploy/docker-compose.infra.yml`.
- `link-api` exposes `GET /healthz`, `GET /readyz`, `GET /{code}` redirect, and management + analytics APIs.
- `link-rpc` exposes readiness, administrator auth, short-link management, redirect resolution, visit recording, and stats query.
- Stage 3 migrations create `admin_user` and `short_link`.
- Stage 5 migrations create `visit_event` and `link_daily_stat`.
- `GET /{code}` returns 302/404/403/410; `AnalyticsMiddleware` fires a non-blocking `RecordVisit` RPC goroutine on every 302.
- `GET /admin/links/{id}/stats` returns daily PV/UV for a short link within a date range.
- Administrator login and profile APIs use JWT-based authentication.
- Authenticated short-link create, list, detail, update, and soft-delete management APIs are implemented.
- Admin SPA under `web/admin/` provides login, link management, and analytics views.
- Local HTTP smoke request assets are available under `tests/httpyac/`.
- `go.mod` is initialized with `github.com/aliaxy/zero-link`.

## Requirements

- Go 1.26 or newer.
- Node.js 22 and pnpm (included in the Nix development shell).
- Docker with Compose support.
- Optional: Nix development shell from `flake.nix`.

If `go`, `node`, or `pnpm` are not available in your normal shell, enter the Nix development environment first:

```bash
nix develop
```

## Local Development

Create local configuration from the committed examples:

```bash
cp .env.example .env.local
cp etc/link-api.example.yaml etc/link-api.local.yaml
cp etc/link-rpc.example.yaml etc/link-rpc.local.yaml
cp web/admin/.env.example web/admin/.env.local
```

The `*.example.*` files are committed templates. The `*.local.*` files are for machine-local values and are ignored by Git.

Start infrastructure:

```bash
make infra-up
make migrate-up
```

Run backend services:

```bash
make run-rpc
make run-api
```

Run the admin UI dev server:

```bash
make web-install
make web-dev
```

Open `http://localhost:5173` and sign in with the seeded administrator (`admin` / `zerolink`).

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
