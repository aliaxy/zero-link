# zero-link

zero-link is a go-zero-oriented short-link system. Stages 1–8 are complete.

## Roadmap

- **CI/CD**: GitHub Actions pipeline for lint, test, and build on pull requests.
- **Containerisation**: Docker Compose application services for one-command local stack startup.
- **Kubernetes**: Helm chart or Kustomize manifests for production deployment.
- **Observability**: Grafana + Prometheus dashboards and alerting rules.
- **Link QR codes**: Generate and serve QR code images for each short link.
- **Multi-tenancy**: Per-workspace isolation of links, stats, and administrators.

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
