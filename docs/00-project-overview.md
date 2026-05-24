# zero-link Project Overview

## Background

zero-link is a short-link system built around go-zero. It should be more than a framework demo: the project needs clear product boundaries, local development support, deployment paths, and enough engineering discipline to evolve safely.

The first product shape is an administrator-managed short-link platform. Administrators create and manage links, public visitors use short links for redirects, and operators can run the system locally or deploy it later through Docker Compose.

## Goals

- Provide a short-link service that can be deployed and operated.
- Keep local development fast by running MySQL and Redis through Docker Compose while Go services run locally.
- Use go-zero API and RPC services to keep HTTP concerns separate from core business logic.
- Preserve a path to future multi-user, tenant, analytics, and SaaS features.
- Make operational behavior visible through logs, metrics, traces, and health checks.

## Users

- **Administrator**: logs into the management UI, creates short links, edits metadata, disables links, and views statistics.
- **Visitor**: opens a short link and expects a fast redirect to the original URL.
- **Operator**: starts dependencies, runs migrations, monitors logs and metrics, and deploys the service.

## Core Scenarios

- Create a short link with a generated code.
- Create a short link with a custom code.
- Redirect a visitor from `/{code}` to the original URL.
- View link lists, link details, and visit statistics.
- Disable or expire a link and ensure it can no longer redirect.
- Recover redirect capability after cache loss by reading from MySQL.

## Non-Goals For The First Stage

- Public user registration.
- Multi-tenant account isolation.
- Billing, plans, and quotas.
- Complex marketing or operations dashboards.
- Custom domains.
- Public Open API tokens.
- Kubernetes deployment.

These are valid future directions, but they should not make the first implementation harder to finish.

## Success Criteria

- MySQL and Redis can be started locally through Docker Compose.
- API and RPC services can run locally against those dependencies.
- An administrator can create, view, update, disable, and inspect short links.
- A visitor can access a short code and receive the correct HTTP redirect or error status.
- Visit events are recorded and aggregated into basic statistics.
- Logs and health checks are sufficient to diagnose common local and deployment failures.
- The project has enough documentation to turn into implementation plans without rediscovering the design.
