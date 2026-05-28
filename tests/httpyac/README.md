# HTTP Smoke Requests

These `.http` files are local smoke-test requests for httpyac.

Use local variables for secrets and generated values. Do not commit real bearer tokens.

Suggested local variables:

```json
{
  "baseUrl": "http://127.0.0.1:8080",
  "adminUsername": "admin",
  "adminPassword": "zerolink",
  "token": "replace-with-login-token",
  "linkId": 1
}
```
