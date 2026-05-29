# API Contracts

This directory contains the go-zero `.api` source file for `link-api`.

`link.api` defines all HTTP routes, request/response types, and middleware annotations. It is the source of truth for code generation — do not handwrite generated files.

## Regenerate

```bash
goctl api go \
  --api services/link-api/api/link.api \
  --dir services/link-api \
  --style gozero
```

Run from the repository root. The `--style gozero` flag matches the existing camelCase file naming convention.
