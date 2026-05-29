# link-rpc

go-zero gRPC service for zero-link. Owns credential verification, business validation, and all MySQL/Redis access.

## RPC Methods

| Method | Description |
|---|---|
| `Check` | Readiness check (MySQL + Redis connectivity) |
| `AuthenticateAdmin` | Verify credentials, return admin profile |
| `GetAdminProfile` | Fetch admin by ID |
| `CreateShortLink` | Create and persist a short link |
| `ListShortLinks` | Paginated list with optional keyword/status filter |
| `GetShortLink` | Fetch single short link by ID |
| `UpdateShortLink` | Update mutable fields, invalidate Redis cache |
| `DeleteShortLink` | Soft-delete, invalidate Redis cache |
| `ResolveShortLink` | Redis-cached code lookup for redirect |
| `RecordVisit` | Insert visit event, upsert daily PV stat |
| `GetLinkStats` | Query daily PV/UV for a date range (max 90 days) |

## Configuration

Copy the example config and edit for local values:

```bash
cp etc/link-rpc.example.yaml etc/link-rpc-local.yaml
```

Key fields: `DataSource`, `CacheRedis`, `Analytics.IPSalt`.

## Run

```bash
make run-rpc
```
