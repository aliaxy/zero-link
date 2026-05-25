# Migrations

This directory contains versioned SQL migrations managed by `golang-migrate`.

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

## Stage 3 Initial Schema

`000001_stage3_initial_schema` creates:

- `admin_user`
- `short_link`

It also seeds a local/dev administrator:

- Username: `admin`
- Password: `zerolink`

The migration stores only the bcrypt password hash. The plaintext password is documented here for local development and verification.

## Policy

- Add one `.up.sql` and one `.down.sql` file for every migration version.
- Review SQL before applying it to shared environments.
- Keep business schema migrations separate from generated go-zero code changes.
