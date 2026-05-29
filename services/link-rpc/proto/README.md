# RPC Contracts

This directory contains protobuf source files for `link-rpc`.

Proto files live under versioned package directories (`link/v1`) and use absolute `go_package` values without explicit Go package aliases.

## Regenerate

Run from the repository root:

```bash
goctl rpc protoc services/link-rpc/proto/link/v1/link.proto \
  --go_out=services/link-rpc/pb \
  --go_opt=paths=source_relative \
  --go-grpc_out=services/link-rpc/pb \
  --go-grpc_opt=paths=source_relative \
  --zrpc_out=services/link-rpc \
  --proto_path=services/link-rpc/proto
```

Accept goctl's generated package names, client directories, and exported client identifiers as the source of truth.
