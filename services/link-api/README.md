# link-api

go-zero HTTP API service for zero-link.

## Endpoints

### Public

| Method | Path | Description |
|---|---|---|
| GET | `/healthz` | Process liveness check |
| GET | `/readyz` | Readiness check (RPC, MySQL, Redis) |
| GET | `/:code` | Short link redirect (302/404/403/410) |
| POST | `/admin/login` | Administrator sign in, returns JWT |

### Authenticated (Bearer token required)

| Method | Path | Description |
|---|---|---|
| GET | `/admin/profile` | Authenticated administrator profile |
| POST | `/admin/links` | Create short link |
| GET | `/admin/links` | List short links (paginated) |
| GET | `/admin/links/:id` | Get short link detail |
| PATCH | `/admin/links/:id` | Update mutable fields |
| DELETE | `/admin/links/:id` | Soft-delete short link |
| GET | `/admin/links/:id/stats` | Daily PV/UV stats for a date range |

## Configuration

Copy the example config and edit for local values:

```bash
cp etc/link-api.example.yaml etc/link-api-local.yaml
```

Key fields: `Auth.Secret`, `Auth.TokenTTLSeconds`, `Cors.AllowOrigins`, `LinkRPC.Target`.

## Run

```bash
make run-api
```
