# zero-link

zero-link is a go-zero-oriented short-link system. The project is currently in a skeleton-only engineering stage: repository structure, local infrastructure, and documentation are prepared, while go-zero service contracts, generated code, and business schema are intentionally deferred.

## Current Stage

Stage 1 builds the project skeleton without business implementation:

- MySQL and Redis are provided by `deploy/docker-compose.infra.yml`.
- Service directories are reserved under `services/`.
- Future migration, API, RPC, and admin UI locations are reserved.
- `go.mod` is initialized with `github.com/aliaxy/zero-link`.

## Requirements

- Go 1.26 or newer.
- Docker with Compose support.
- Optional: Nix development shell from `flake.nix`.

If `go` is not available in your normal shell, enter the Nix development environment first.

## Local Development

```bash
make infra-up
```

Generation, migrations, service run commands, and business implementation are reserved for later stages. The current repository contains skeleton directories, infra config, and documentation only.

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
