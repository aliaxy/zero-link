# web/admin

Vue 3 + Vite + Element Plus admin SPA for zero-link.

## Stack

- **Vue 3** — Composition API
- **Vite 6** — dev server and build
- **Element Plus** — UI component library
- **Pinia** — state management (auth store)
- **Vue Router 4** — client-side routing
- **ECharts + vue-echarts** — PV/UV analytics chart
- **axios** — HTTP client with auth interceptor
- **dayjs** — date formatting
- **Geist Variable** — typography

## Pages

| Route | Description |
|---|---|
| `/login` | Administrator sign in |
| `/links` | Link list with search, filter, create, edit, delete |
| `/links/:id` | Link detail with edit form and short URL preview |
| `/links/:id/stats` | Daily PV/UV chart and data table |

## Local Development

Install dependencies (requires pnpm from the Nix shell):

```bash
make web-install
```

Create a local environment file:

```bash
cp web/admin/.env.example web/admin/.env.local
```

Start the dev server (proxies `/api/*` and `/:code` to link-api):

```bash
make web-dev
```

Open `http://localhost:5173`.

## Build

```bash
make web-build
```

Output lands in `web/admin/dist/` (ignored by Git).

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `VITE_API_BASE_URL` | `/api` | Prefix for API calls (Vite proxy strips it) |
| `VITE_API_TARGET` | `http://127.0.0.1:8080` | link-api address for Vite proxy |
| `VITE_SHORT_LINK_BASE` | `window.location.origin` | Base URL shown in short link preview |
