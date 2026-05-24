# System Architecture

## Overview

zero-link uses a go-zero API + RPC architecture. The HTTP service handles web-facing concerns, while the RPC service owns short-link business behavior and persistence.

```text
Visitor / Admin Browser
        |
        v
     link-api
  HTTP, auth, validation,
  redirect route, UI assets
        |
        v
     link-rpc
  business rules, cache,
  storage, analytics
    |           |
    v           v
  MySQL       Redis
```

## Components

### link-api

Responsibilities:

- Expose management HTTP APIs.
- Expose `GET /{code}` redirect route.
- Serve or proxy the simple management UI.
- Validate HTTP request shape.
- Apply administrator authentication middleware.
- Translate domain errors into HTTP responses.
- Call `link-rpc` for business operations.

The API service should not contain storage-specific business logic.

### link-rpc

Responsibilities:

- Create, update, fetch, and resolve short links.
- Generate and validate short codes.
- Read and write MySQL records.
- Manage Redis cache and cache invalidation.
- Record visit events.
- Query aggregated statistics.

The RPC service is the primary boundary for business rules.

### MySQL

Responsibilities:

- Store administrator accounts.
- Store short-link source-of-truth records.
- Store visit events.
- Store aggregated statistics.

### Redis

Responsibilities:

- Cache short-code resolution data.
- Hold rate-limit counters.
- Optionally support UV de-duplication.
- Optionally buffer visit events in a later iteration.

### Management UI

Responsibilities:

- Provide administrator-facing workflows.
- Call management HTTP APIs.
- Show validation and operation errors clearly.

The first UI should stay small and can be served by `link-api`.

## Local Development Architecture

- `deploy/docker-compose.infra.yml` starts MySQL and Redis.
- `link-api` and `link-rpc` run on the host machine.
- Local config points to `127.0.0.1` ports exposed by Compose.
- This keeps logs, debugger attachment, and rebuild cycles simple.

## Deployment Architecture

After core behavior is stable:

- Add Dockerfiles for `link-api` and `link-rpc`.
- Add a full `docker-compose.yml` containing MySQL, Redis, RPC, API, and UI.
- Compose-internal config uses service names such as `mysql`, `redis`, and `link-rpc`.
- MySQL data uses a named volume.

## Service Discovery

The first stage should not require etcd. Local development can use direct RPC targets to reduce moving parts.

When the system needs multiple RPC instances or dynamic discovery, introduce the go-zero standard service discovery path with etcd and document the migration.

## Dependency Direction

- HTTP handlers depend on API logic.
- API logic depends on RPC clients.
- RPC logic depends on repositories, cache clients, and analytics components.
- Repositories do not depend on HTTP or RPC transport types.

This keeps business code testable without a running HTTP server.
