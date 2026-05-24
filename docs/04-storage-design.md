# Storage Design

## Principles

- MySQL is the source of truth.
- Redis is a cache and coordination aid, not durable storage.
- Use explicit SQL or go-zero generated models; do not introduce an ORM.
- Pass `context.Context` through all database operations.
- Use parameterized queries for all user-controlled values.
- Convert `sql.ErrNoRows` into domain-level not-found errors.

## Tables

### admin_user

Purpose: store administrator accounts.

Columns:

- `id BIGINT PRIMARY KEY AUTO_INCREMENT`
- `username VARCHAR(64) NOT NULL`
- `password_hash VARCHAR(255) NOT NULL`
- `status TINYINT NOT NULL`
- `created_at DATETIME NOT NULL`
- `updated_at DATETIME NOT NULL`

Indexes:

- Unique index on `username`.
- Index on `status`.

### short_link

Purpose: store short-code mappings.

Columns:

- `id BIGINT PRIMARY KEY AUTO_INCREMENT`
- `code VARCHAR(32) NOT NULL`
- `origin_url TEXT NOT NULL`
- `title VARCHAR(255) NOT NULL DEFAULT ''`
- `description VARCHAR(1024) NOT NULL DEFAULT ''`
- `status TINYINT NOT NULL`
- `expire_at DATETIME NULL`
- `created_by BIGINT NOT NULL`
- `created_at DATETIME NOT NULL`
- `updated_at DATETIME NOT NULL`
- `deleted_at DATETIME NULL`

Indexes:

- Unique index on `code`.
- Index on `status`.
- Index on `created_at`.
- Index on `created_by`.
- Index on `expire_at`.

### link_visit_event

Purpose: store visit-level analytics events.

Columns:

- `id BIGINT PRIMARY KEY AUTO_INCREMENT`
- `link_id BIGINT NOT NULL`
- `code VARCHAR(32) NOT NULL`
- `visited_at DATETIME NOT NULL`
- `ip_hash VARCHAR(128) NOT NULL DEFAULT ''`
- `user_agent TEXT NULL`
- `referer TEXT NULL`
- `device VARCHAR(64) NOT NULL DEFAULT ''`
- `browser VARCHAR(64) NOT NULL DEFAULT ''`
- `os VARCHAR(64) NOT NULL DEFAULT ''`
- `country VARCHAR(64) NOT NULL DEFAULT ''`
- `province VARCHAR(64) NOT NULL DEFAULT ''`
- `city VARCHAR(64) NOT NULL DEFAULT ''`

Indexes:

- Index on `link_id, visited_at`.
- Index on `code, visited_at`.
- Index on `visited_at`.

### link_daily_stat

Purpose: store daily aggregate analytics.

Columns:

- `id BIGINT PRIMARY KEY AUTO_INCREMENT`
- `link_id BIGINT NOT NULL`
- `stat_date DATE NOT NULL`
- `pv BIGINT NOT NULL DEFAULT 0`
- `uv BIGINT NOT NULL DEFAULT 0`
- `created_at DATETIME NOT NULL`
- `updated_at DATETIME NOT NULL`

Indexes:

- Unique index on `link_id, stat_date`.
- Index on `stat_date`.

## Consistency Strategy

### Create Short Link

1. Validate input.
2. Generate or validate code.
3. Insert into MySQL.
4. Do not require immediate Redis write; redirect can populate cache on first read.

### Resolve Short Link

1. Try Redis.
2. On miss, read MySQL.
3. Validate status and expiration.
4. Backfill Redis when redirectable.

### Update Short Link

1. Update MySQL.
2. Delete the Redis cache key for the code.
3. Let later redirects repopulate cache.

### Record Visit

1. Redirect path emits a lightweight visit event.
2. Event storage failure is logged and counted but does not cancel a valid redirect.
3. Aggregation reads visit events and upserts daily stats.

## Migration Policy

- Schema changes must be versioned.
- Migration SQL should be reviewed by humans.
- Rollback scripts are required for destructive changes.
- Local development must support recreating a clean schema.
