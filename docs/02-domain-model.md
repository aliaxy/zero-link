# Domain Model

## Core Objects

### AdminUser

Represents an administrator who can access management features.

Fields:

- `id`: unique identifier.
- `username`: login name.
- `password_hash`: password hash produced by bcrypt or Argon2id.
- `status`: active or disabled.
- `created_at`: creation time.
- `updated_at`: last update time.

### ShortLink

Represents a short-code mapping to an original URL.

Fields:

- `id`: unique identifier.
- `code`: public short code.
- `origin_url`: target URL.
- `title`: display title.
- `description`: optional note.
- `status`: lifecycle state.
- `expire_at`: optional expiration time.
- `created_by`: administrator ID.
- `created_at`: creation time.
- `updated_at`: last update time.
- `deleted_at`: soft delete time.

### VisitEvent

Represents one redirect attempt worth recording.

Fields:

- `id`: unique identifier.
- `link_id`: short-link ID.
- `code`: short code at visit time.
- `visited_at`: visit timestamp.
- `ip_hash`: hashed or anonymized visitor IP.
- `user_agent`: raw User-Agent, subject to retention policy.
- `referer`: HTTP referer.
- `device`: parsed device category.
- `browser`: parsed browser.
- `os`: parsed operating system.
- `country`, `province`, `city`: optional geographic fields.

### DailyStat

Represents aggregated daily analytics.

Fields:

- `id`: unique identifier.
- `link_id`: short-link ID.
- `stat_date`: aggregation date.
- `pv`: page views.
- `uv`: unique visitors.
- `created_at`: creation time.
- `updated_at`: last update time.

## Short-Link Status

- `unknown`: invalid zero-value state used defensively in code.
- `active`: redirects are allowed if the link is not expired.
- `disabled`: redirects are blocked by administrator action.
- `expired`: redirects are blocked because expiration time has passed.
- `deleted`: hidden from normal use and treated as unavailable publicly.

The database may store `active`, `disabled`, and `deleted`; `expired` can be derived from `expire_at`.

## Lifecycle

1. **Created**: a link is inserted with a unique code and valid original URL.
2. **Active**: the link can redirect while enabled and unexpired.
3. **Disabled**: the administrator blocks redirects.
4. **Expired**: the configured expiration time has passed.
5. **Archived**: the administrator soft deletes the link.

## Domain Rules

- A short code is immutable after creation.
- MySQL is authoritative for link existence, ownership, status, and URL.
- Redis only accelerates resolution and can be rebuilt from MySQL.
- Updating the original URL, status, expiration, or deletion state invalidates cached resolution data.
- Analytics are eventually consistent and cannot prevent a valid redirect.
- Soft-deleted links must not be reassigned automatically.
