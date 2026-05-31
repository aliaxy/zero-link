# Observability Design

## Goals

- Diagnose redirect failures quickly.
- Understand dependency health.
- Measure cache effectiveness and latency.
- Detect analytics write failures.
- Avoid duplicate error logs and noisy high-cardinality labels.

## Logging

Use structured logs via `logx.Infow` / `logx.Errorw`. The `caller=file:line` field is added automatically by go-zero logx — do not add explicit `func` fields.

Important request fields:

- request ID.
- method.
- route pattern.
- status code.
- duration.
- administrator ID when authenticated.
- short code for redirect requests.

Error logging rule:

- Internal functions return wrapped errors.
- HTTP or RPC boundaries log errors once.
- Do not both log and return the same error at each layer.

Implemented structured log points:

| Layer | Location | Event | Fields |
|-------|----------|-------|--------|
| link-api | `loginlogic.go` | login attempt | username |
| link-api | `loginlogic.go` | login failed | username, reason |
| link-api | `loginlogic.go` | login success | admin_id, username |
| link-api | `authmiddleware.go` | admin authenticated | admin_id, username |
| link-api | `redirectlogic.go` | redirect rpc failed | code, error |
| link-api | `redirectlogic.go` | redirect resolved | code, url |
| link-rpc | `resolveshortlinklogic.go` | resolve hit/miss/disabled/expired | code, result |
| link-rpc | `resolveshortlinklogic.go` | resolve failed | code, error |
| link-rpc | `authenticateadminlogic.go` | authenticate failed | username |
| link-rpc | `authenticateadminlogic.go` | authenticate success | admin_id, username |
| link-rpc | `recordvisitlogic.go` | upsert daily stat failed | link_id, code, stat_date, error |

## Metrics

Implemented business metrics (in `services/link-rpc/internal/metrics/`):

- `zerolink_redirect_requests_total{result}` — result values: hit, miss, disabled, expired, error.
- `zerolink_analytics_events_total{result}` — result values: success, error.

Prometheus HTTP endpoints:

- `link-api`: `http://127.0.0.1:9100/metrics` (go-zero framework metrics).
- `link-rpc`: `http://127.0.0.1:9101/metrics` (framework + business metrics). `Middlewares.Prometheus: true` enables built-in gRPC RPC latency instrumentation.

Recommended additional metrics (not yet implemented):

- `http_request_duration_seconds{method,route,status}`
- `redirect_duration_seconds{result}`
- `redirect_cache_hits_total` / `redirect_cache_misses_total`
- `mysql_errors_total{operation}`
- `redis_errors_total{operation}`

Avoid labels containing raw URL, IP, user ID, or arbitrary short-code values.

## Tracing

Trace path:

```text
HTTP request -> link-api handler -> link-rpc method -> Redis/MySQL -> response
```

Trace spans should include:

- Route pattern.
- RPC method.
- Cache hit or miss.
- Database operation name.
- Error status when failures occur.

## Health Checks

### healthz

Confirms the process is alive and can respond. It should not depend on MySQL or Redis.

### readyz

Checks whether the service can handle real traffic.

For `link-api`:

- RPC connectivity.

For `link-rpc`:

- MySQL ping.
- Redis ping.

## Troubleshooting Scenarios

### Redirect Is Slow

Check:

- Redirect latency histogram.
- Redis hit ratio.
- MySQL query latency.
- RPC latency.

### Redis Hit Rate Is Low

Check:

- Cache TTL.
- Cache invalidation frequency.
- Redis errors.
- Whether redirects are using many one-off codes.

### Analytics Are Delayed

Check:

- Analytics event failure metric.
- Worker queue depth if a queue exists.
- MySQL insert errors.
- Aggregation job duration.

### Admin Create Fails

Check:

- Management request log.
- Validation error code.
- MySQL unique constraint errors.
- RPC create method latency and error rate.
