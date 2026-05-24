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
