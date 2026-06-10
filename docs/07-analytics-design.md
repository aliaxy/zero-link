# Analytics Design

## Definitions

- **PV**: every recorded visit event for a link.
- **UV**: unique visitor count for a link within a time window, based on `ip_hash` plus optional User-Agent fingerprint.
- **Referer**: source page from the HTTP `Referer` header.
- **Device**: parsed category such as desktop, mobile, tablet, bot, or unknown.

## Visit Event Fields

Each event should capture:

- Link ID.
- Code.
- Visit time.
- IP hash.
- User-Agent.
- Referer.
- Device.
- Browser.
- Operating system.
- Optional geographic fields.

## User-Agent Parsing

The redirect path can store the raw User-Agent and defer parsing to an asynchronous process. This avoids adding CPU-heavy parsing to the critical path.

If parsing synchronously in the first implementation, failures must produce `unknown` values and never block redirect.

## Referer Parsing

Referer should be normalized into a domain for aggregate display. Raw referer may be retained initially, but retention and privacy policy should be revisited before production use.

## IP Privacy

Store `ip_hash` instead of raw IP by default.

Rules:

- Use a server-side salt from configuration.
- Do not log the salt.
- Rotate salt only with awareness that UV continuity will change.
- Treat IP-derived data as sensitive.

## Write Strategy

Current implementation:

- `AnalyticsMiddleware` dispatches a job to a bounded channel-backed worker pool after every 302 response.
- Pool size: 8 workers, 2 000-slot buffer. Jobs are dropped (with a log line) when the buffer is full,
  preventing memory exhaustion under traffic spikes.
- Each worker opens a 3-second context for the `RecordVisit` RPC call.
- Redirect response is not delayed by analytics write completion.
- Failures increment metrics and emit structured logs.

Future stage:

- Move event buffering to Redis Stream, Kafka, or another queue.
- Run a dedicated analytics consumer.
- Use OLAP storage if MySQL aggregation becomes too expensive.

## Aggregation Strategy

- Aggregate into `link_daily_stat`.
- Upsert by `link_id + stat_date`.
- Store PV and UV.
- UI reads aggregates for dashboards.
- UI should not scan large event tables for routine summary views.

## Delay Tolerance

Analytics can be delayed by seconds or minutes. The UI should avoid promising real-time precision unless the implementation supports it.

## Failure Policy

- Successful redirects remain successful if analytics fails.
- Analytics failures are logged once at the boundary and counted in metrics.
- Repeated analytics failures should be visible in readiness or alerting only if they threaten data integrity requirements.
