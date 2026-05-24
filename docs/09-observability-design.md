# Observability Design

## Goals

- Diagnose redirect failures quickly.
- Understand dependency health.
- Measure cache effectiveness and latency.
- Detect analytics write failures.
- Avoid duplicate error logs and noisy high-cardinality labels.

## Logging

Use structured logs.

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

## Metrics

Recommended metrics:

- `http_requests_total{method,route,status}`
- `http_request_duration_seconds{method,route,status}`
- `redirect_requests_total{result}`
- `redirect_duration_seconds{result}`
- `redirect_cache_hits_total`
- `redirect_cache_misses_total`
- `analytics_events_total{result}`
- `rpc_requests_total{method,status}`
- `rpc_request_duration_seconds{method,status}`
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
