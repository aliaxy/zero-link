# Storage Design

## Principles

- MySQL is the source of truth.
- Redis is a cache and coordination aid, not durable storage.
- Use explicit SQL or go-zero generated models; do not introduce an ORM.
- Pass `context.Context` through all database operations.
- Use parameterized queries for all user-controlled values.
- Convert `sql.ErrNoRows` into domain-level not-found errors.
- Use `golang-migrate` with reviewed versioned SQL for schema changes.

## Stage 3 Tables

Stage 3 introduces the first business schema for administrator authentication and short-link management only.

### admin_user

Purpose: store administrator accounts.

Columns:

- `id BIGINT PRIMARY KEY AUTO_INCREMENT`
- `username VARCHAR(64) NOT NULL`
- `password_hash VARCHAR(255) NOT NULL`
- `status TINYINT NOT NULL`
- `created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP`
- `updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`

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
- `created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP`
- `updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`
- `deleted_at DATETIME NULL`

Indexes:

- Unique index on `code`.
- Index on `status`.
- Index on `created_at`.
- Index on `created_by`.
- Index on `expire_at`.

Constraints:

- Foreign key from `short_link.created_by` to `admin_user.id`.

## Stage 5 Analytics Tables

Stage 5 adds visit event and daily stat tables under `migrations/000002_stage5_analytics`.

### visit_event

Purpose: store visit-level analytics events.

Columns:

- `id BIGINT PRIMARY KEY AUTO_INCREMENT`
- `link_id BIGINT NOT NULL`
- `code VARCHAR(32) NOT NULL`
- `visited_at DATETIME(6) NOT NULL`
- `ip_hash VARCHAR(64) NOT NULL DEFAULT ''`
- `user_agent VARCHAR(512) NOT NULL DEFAULT ''`
- `referer VARCHAR(1024) NOT NULL DEFAULT ''`
- `device VARCHAR(16) NOT NULL DEFAULT 'unknown'`

Indexes:

- Index on `link_id, visited_at`.

Notes:

- No foreign key from `link_id` to `short_link.id`; analytics writes are fire-and-forget and must not fail if a link is deleted.
- `ip_hash` is HMAC-SHA256 of the raw IP with a server-side salt stored only in `link-rpc` config.
- `device` is one of `bot`, `mobile`, `desktop`, `unknown` (detected from User-Agent).
- `browser`, `os`, `country`, `province`, and `city` are deferred to a later stage.

### link_daily_stat

Purpose: store daily aggregate analytics.

Columns:

- `id BIGINT PRIMARY KEY AUTO_INCREMENT`
- `link_id BIGINT NOT NULL`
- `stat_date DATE NOT NULL`
- `pv BIGINT NOT NULL DEFAULT 0`
- `uv BIGINT NOT NULL DEFAULT 0`

Indexes:

- Unique index on `link_id, stat_date` (`uk_link_daily_stat`).

Notes:

- `UpsertPV` uses `INSERT ... ON DUPLICATE KEY UPDATE pv = pv + 1`.
- `uv` is set to 1 on first insert and not updated on subsequent visits (Stage 5 approximation; true UV counting is deferred).
- No foreign key from `link_id`; same reason as `visit_event`.

## Consistency Strategy

## Stage 8 Data Lifecycle Tables

Stage 8 adds two tables under `migrations/000003_data_retention` to support data retention and code reuse prevention.

### short_link_archive

Purpose: permanent archive of soft-deleted short-link rows after they exceed the retention window (default 365 days).

Columns:

- `id BIGINT NOT NULL` — preserves the original `short_link.id`; not auto-increment.
- `code VARCHAR(32) NOT NULL`
- `origin_url TEXT NOT NULL`
- `title VARCHAR(255) NOT NULL DEFAULT ''`
- `description VARCHAR(1024) NOT NULL DEFAULT ''`
- `status TINYINT NOT NULL`
- `expire_at DATETIME NULL`
- `created_by BIGINT NOT NULL`
- `created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP`
- `updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`
- `deleted_at DATETIME NULL`

Indexes:

- `KEY idx_sla_code (code)`.
- `KEY idx_sla_deleted_at (deleted_at)`.

Notes:

- No `UNIQUE KEY` on `code` and no foreign keys. `INSERT IGNORE` ensures idempotency on crash-restart.
- Preserves original identity: `id` equals the original `short_link.id`.

### reserved_code

Purpose: permanently reserve codes after archival so they can never be reassigned.

Columns:

- `code VARCHAR(32) NOT NULL` — primary key.
- `reserved_at DATETIME NOT NULL`.

Notes:

- Primary key on `code` provides O(1) lookup.
- Written during the archival step: after a row is moved from `short_link` to `short_link_archive`, its code is inserted here with `INSERT IGNORE`.
- `CreateShortLink` checks this table for custom codes before allowing creation.

## Data Retention Policy

The cleanup runner in `services/link-rpc/internal/cleanup/` executes once on startup and then every 24 hours. Retention windows and batch sizes are controlled by `Retention.*` config fields (with zero-value defaults applied in `NewServiceContext`):

| Data | Default retention | Cleanup action |
|------|-------------------|----------------|
| `visit_event` | 90 days | Batch DELETE by `visited_at` |
| `short_link` (soft-deleted) | 365 days | INSERT IGNORE to archive → INSERT IGNORE to reserved_code → DELETE |
| `link_daily_stat` | 730 days | Batch DELETE by `stat_date` |

Batch deletes use `LIMIT CleanupBatchSize` (default 1000) in a loop to avoid long-running table locks. The archival step is idempotent: crash-restart safety is provided by `INSERT IGNORE` at each step rather than a transaction.

## Consistency Strategy

### Create Short Link

1. Validate input.
2. Generate or validate code.
3. Insert into MySQL.
4. Do not require an immediate Redis write in Stage 3.

### Resolve Short Link

Redirect resolution is deferred until Stage 4. The intended future flow is:

1. Try Redis.
2. On miss, read MySQL.
3. Validate status and expiration.
4. Backfill Redis when redirectable.

### Update Short Link

1. Update MySQL.
2. Keep cache invalidation requirements documented for Stage 4.
3. Do not add Redis redirect cache behavior during Stage 3.

### Record Visit

1. `AnalyticsMiddleware` in `link-api` fires a goroutine after every 302 redirect.
2. The goroutine calls `RecordVisit` RPC with a 3-second timeout context.
3. `RecordVisit` looks up `link_id` via `FindOneByCode` (Redis-cached), hashes the IP, inserts `visit_event`, then upserts `link_daily_stat`.
4. `UpsertPV` failure is logged but does not propagate; the redirect is unaffected.

## Migration Policy

- Schema changes must be versioned.
- Use `golang-migrate` for local and future deployment migrations.
- Migration SQL should be reviewed by humans.
- Rollback scripts are required for destructive changes.
- Local development must support recreating a clean schema.
- Stage 3 should create `admin_user` and `short_link` before any generated models or business logic depend on them.
- Stage 5 should create `visit_event` and `link_daily_stat` before analytics models or logic depend on them.
- Local development may seed one default administrator through migration SQL so login can be verified without manual setup.
- Seeded administrator passwords must be stored as bcrypt hashes; plaintext passwords must not be committed outside documented local examples.
