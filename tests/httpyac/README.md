# HTTP Smoke Requests

These `.http` files are local smoke-test requests for httpyac.

Use `http-client.example.env.json` as the committed variable reference. Keep real passwords,
bearer tokens, and generated IDs in `http-client.private.env.json`, which is ignored by Git.

Create your local private environment file:

```bash
cp tests/httpyac/http-client.example.env.json tests/httpyac/http-client.private.env.json
```

Suggested local variables use the `local` environment:

```json
{
  "local": {
    "baseUrl": "http://127.0.0.1:8080",
    "adminUsername": "admin",
    "adminPassword": "zerolink",
    "token": "replace-with-login-token",
    "linkId": 1,
    "linkCode": "replace-with-short-code"
  }
}
```

Run a named request:

```bash
httpyac send tests/httpyac/admin.http -e local -n login
```

Run every request in a file:

```bash
httpyac send tests/httpyac/health.http -e local --all
httpyac send tests/httpyac/admin.http -e local --all
httpyac send tests/httpyac/analytics.http -e local --all
```

After login, copy the returned token into `http-client.private.env.json` before running authenticated
management requests. After creating a link, copy the returned `id` into `linkId` and the returned
`code` into `linkCode` before running detail, update, delete, redirect, or stats requests.
