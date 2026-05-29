# Migrations

Versioned SQL migrations managed by `golang-migrate`.

## Local Usage

Create local environment values first:

```bash
cp .env.example .env.local
```

Ensure `.env.local` contains `ZERO_LINK_MIGRATE_DSN`, then run:

```bash
make infra-up
make migrate-up
```

Roll back the latest local migration step with:

```bash
make migrate-down
```

## Versions

### 000001 — Stage 3 Initial Schema

Creates `admin_user` and `short_link`.

Seeds a local/dev administrator:

- Username: `admin`
- Password: `zerolink`

The migration stores only the bcrypt password hash. The plaintext password is documented here for local development only.

### 000002 — Stage 5 Analytics

Creates `visit_event` and `link_daily_stat`.

## Policy

- Add one `.up.sql` and one `.down.sql` file for every migration version.
- Review SQL before applying it to shared environments.
- Keep business schema migrations separate from generated go-zero code changes.
