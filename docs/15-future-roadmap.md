# Future Roadmap

## Account And Organization Features

- Multiple administrators.
- Role-based permissions.
- Teams or organizations.
- Multi-tenant data isolation.
- Audit logs for administrator actions.

## Link Features

- Custom domains.
- Batch short-link creation.
- QR code generation.
- Link tags and folders.
- Link templates.
- Password-protected links.
- Geo or device-based redirect rules.

## Public API

- API tokens.
- Token scopes.
- Per-token rate limits.
- OpenAPI documentation.
- SDK examples.

## Analytics Evolution

- Redis Stream or message queue for visit events.
- Dedicated analytics consumer service.
- ClickHouse or another OLAP store for high-volume analytics.
- More dimensions: browser, OS, country, city, campaign tags.
- Export to CSV.

## Security And Abuse Control

- Target domain allowlist or denylist.
- Malware or phishing checks.
- Link review workflow.
- Bot filtering.
- Abuse reports.
- Admin MFA.

## Commercialization

- Plans and quotas.
- Usage metering.
- Billing integration.
- Custom domain limits.
- Team member limits.

## Operations

- CI/CD.
- Automated vulnerability scanning.
- Grafana dashboards based on Stage 7 foundational metrics (`zerolink_redirect_requests_total`, `zerolink_analytics_events_total`) and go-zero built-in HTTP/gRPC latency metrics.
- Alert rules (e.g., redirect error rate, analytics failure rate, rate-limit saturation).
- Backup and restore playbooks.
- Kubernetes deployment.
- Blue-green or rolling deployments.

## Architecture Evolution

- Introduce etcd for go-zero service discovery when multi-instance RPC deployment is needed.
- Split analytics into a separate service when event volume grows.
- Add a dedicated frontend deployment if the management UI grows beyond a simple embedded admin.
- Add read replicas or caching improvements if MySQL read load becomes a bottleneck.
