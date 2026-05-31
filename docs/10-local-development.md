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
- `web/admin/.env.example`: committed frontend environment variable example.
- `web/admin/.env.local`: machine-local frontend environment values ignored by Git.
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
make web-install
make web-dev
make web-build
make web-preview
```

`make run-rpc` starts the RPC service with `etc/link-rpc-local.yaml`.
`make run-api` starts the API service with `etc/link-api-local.yaml`.
`make web-install` installs frontend dependencies via pnpm.
`make web-dev` starts the Vite dev server at `http://localhost:5173`.
`make web-build` produces a production build under `web/admin/dist/`.
`make web-preview` serves the production build locally for verification.

## Migration Workflow

Uses a `golang-migrate` workflow.

Intended local behavior:

- `make migrate-up` applies all pending migrations to the local MySQL database.
- `make migrate-down` rolls back the latest local migration step.
- Migration files live under `migrations/` and use versioned up/down SQL.
- `000001`: creates `admin_user` and `short_link` (Stage 3).
- `000002`: creates `visit_event` and `link_daily_stat` (Stage 5).
- A local/dev seed administrator may be inserted by migration SQL so login can be tested without manual database edits.

Do not run migrations before `make infra-up` has started MySQL.

## Local Verification Flow

### Backend only

1. Start MySQL and Redis with `make infra-up`.
2. Run `make migrate-up` to apply the Stage 3 and Stage 5 schema.
3. Start `link-rpc` locally with `make run-rpc`.
4. Start `link-api` locally with `make run-api`.
5. Open `http://127.0.0.1:8080/healthz`.
6. Open `http://127.0.0.1:8080/readyz`.
7. Confirm `healthz` succeeds when the API process is alive.
8. Confirm `readyz` succeeds only when API, RPC, MySQL, and Redis are reachable.
9. Log in with the seeded local administrator and create a short link through the management API.
10. Access `GET /{code}` and confirm a `302` redirect; check `visit_event` for a recorded row.
11. Call `GET /admin/links/{id}/stats` and confirm `pv: 1` for today.
12. Open `http://127.0.0.1:9100/metrics` to verify link-api Prometheus metrics are exported.
13. Open `http://127.0.0.1:9101/metrics` and confirm `zerolink_redirect_requests_total` and `zerolink_analytics_events_total` are present.

### Full stack with admin UI

1. Complete steps 1–4 above.
2. Create frontend environment file: `cp web/admin/.env.example web/admin/.env.local`.
3. Install frontend dependencies: `make web-install`.
4. Start the Vite dev server: `make web-dev`.
5. Open `http://localhost:5173` and log in with the seeded administrator.
6. Create a short link; confirm it appears in the list.
7. Click a table row to open link detail; edit and save a change.
8. Open the short URL `http://localhost:5173/:code` in a browser; confirm a redirect.
9. Navigate to the stats view and confirm PV increments after the redirect.

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

## Frontend Environment

Create a local frontend environment file before running `make web-dev`:

```bash
cp web/admin/.env.example web/admin/.env.local
```

Key variables:

- `VITE_API_BASE_URL`: prefix for API calls; defaults to `/api` which the Vite proxy rewrites to the link-api root.
- `VITE_API_TARGET`: link-api address for the Vite proxy; defaults to `http://127.0.0.1:8080`.
- `VITE_SHORT_LINK_BASE`: base URL shown in the link detail short URL preview; defaults to `window.location.origin`.

The Vite dev server proxies two path classes to link-api:

- `/api/*` → strips `/api` prefix and forwards to link-api.
- `/:code` → forwards short link codes directly; the link-api redirect handler returns 302.

Do not commit `web/admin/.env.local`. The `web/admin/dist/` build output is also ignored by Git.

## Tooling Note

The repository includes a Nix development shell with Go, goctl, protobuf tooling, Node.js, and pnpm. If any of these tools are not available in a normal shell, enter the Nix development environment before running local commands:

```bash
nix develop
```
