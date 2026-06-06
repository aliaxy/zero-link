# Deployment Design

## Deployment Stages

### Stage 1: Local-First Development

- MySQL and Redis run through `deploy/docker-compose.infra.yml`.
- `link-api` and `link-rpc` run as local Go processes.
- This is the default development and test workflow.

### Stage 2: Server Process Deployment

- MySQL and Redis can still run through Compose or managed services.
- Go services run as system processes or supervised services.
- Config points to the chosen dependency endpoints.

### Stage 3: Full Docker Compose Deployment

- MySQL.
- Redis.
- `link-rpc`.
- `link-api`.
- Management UI assets or UI service.

This stage is added after core behavior is stable locally.

## Full Compose Layout

Services:

- `mysql`
- `redis`
- `link-rpc`
- `link-api`
- optional `adminer`

Networking:

- API exposes the public HTTP port.
- RPC is internal to the Compose network.
- MySQL and Redis are internal unless local debugging requires exposed ports.
- Prometheus metrics endpoints (port 9100 for link-api, port 9101 for link-rpc) should be reachable by a Prometheus scraper or remain internal; do not expose them publicly.
- Structured logs from both services are JSON-formatted (go-zero logx); route to a log aggregator in production deployments.

## Configuration Differences

Local development:

- MySQL host: `127.0.0.1`
- Redis host: `127.0.0.1`
- RPC target: local host and port.

Compose deployment:

- MySQL host: `mysql`
- Redis host: `redis`
- RPC target: `link-rpc:{port}` or service discovery target.

## Data Persistence

- MySQL uses a named volume.
- Redis persistence is optional in the first stage because Redis is a cache.
- If Redis is later used as an event buffer, persistence requirements must be revisited.

## Startup Order

1. MySQL.
2. Redis.
3. Migrations.
4. `link-rpc`.
5. `link-api`.
6. Management UI if separate.

Readiness checks, not startup order alone, determine whether traffic can be served.

## Rollback Considerations

- Application rollback is straightforward when schema remains backward-compatible.
- Destructive database migrations require explicit rollback scripts and review.
- Cache can be cleared safely because MySQL is authoritative.

## Stage 9 Implementation

Stage 9 delivers the full Docker Compose stack. The files are:

- `services/link-api/Dockerfile` — multi-stage build for link-api
- `services/link-rpc/Dockerfile` — multi-stage build for link-rpc
- `etc/link-api.compose.yaml` — Compose-specific service config
- `etc/link-rpc.compose.yaml` — Compose-specific service config
- `.env.compose.example` — all env vars the Compose stack consumes
- `deploy/docker-compose.yml` — full application and infrastructure Compose file

### Service Topology

Startup dependency chain enforced by `depends_on` conditions:

```
mysql (healthy) ──┐
                  ├── migrate (completed) ──┐
redis (healthy) ──┘                        ├── link-rpc (healthy) ── link-api
                  └───────────────────────┘
```

### Configuration Approach

go-zero's `conf.MustLoad` calls `os.ExpandEnv` on the config file at startup, which expands `${VAR}` references from the container environment. The Compose config files (`etc/link-api.compose.yaml`, `etc/link-rpc.compose.yaml`) are baked into the Docker images using `${VAR}` placeholders.

The `docker-compose.yml` `environment:` block resolves defaults with Docker Compose's `${VAR:-default}` syntax before passing values to containers. This cleanly separates default resolution (Compose) from env var expansion (go-zero).

Secrets (`LINK_API_AUTH_SECRET`, `LINK_RPC_IP_SALT`) have no defaults and must be set in `.env`.

### Operator Quick-Start

```bash
cp .env.compose.example .env.compose.local
# Edit .env.compose.local: set LINK_API_AUTH_SECRET (>=32 chars) and LINK_RPC_IP_SALT
make compose-up
make compose-logs
curl http://localhost:8080/healthz
```

### Port Exposure

Only port 8080 (link-api HTTP) is published to the host. All other ports (9090 gRPC, 9100/9101 metrics, 3306 MySQL, 6379 Redis) remain internal to the Compose network.

### Migration Service

The `migrate` service uses `ghcr.io/golang-migrate/migrate` as a one-shot container. It depends on MySQL being healthy and exits with code 0 on success. `link-rpc` uses `condition: service_completed_successfully` so it only starts after migrations have applied. This requires Docker Compose v2 (already in use).

### MySQL DSN and Timezone

Compose configs use `loc=UTC` in the MySQL DSN instead of `loc=Local`. This avoids any dependency on host timezone configuration and is consistent across containers.

### Mutual Exclusivity

`deploy/docker-compose.yml` and `deploy/docker-compose.infra.yml` are mutually exclusive. Both define the same MySQL and Redis container names and volume names. Running both simultaneously is not supported. Use `docker-compose.infra.yml` for local Go process development and `docker-compose.yml` for the full containerised stack.

### Admin UI

The Vue 3 admin SPA (`web/admin/`) is out of scope for Stage 9. It continues to be served by the Vite dev server during local development and is not packaged into the Compose stack.
