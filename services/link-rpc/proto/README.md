# RPC Contracts

This directory contains protobuf files.

The current protobuf contract includes only a readiness RPC contract. Stage 3 database foundation exists, but business RPC methods are reserved for the next Stage 3 implementation pass.

Proto files live under versioned package directories, such as `link/v1`, and should use absolute `go_package` values without explicit Go package aliases.

RPC generation should be run from the repository root:

```bash
goctl rpc protoc services/link-rpc/proto/link/v1/link.proto \
  --go_out=services/link-rpc/pb \
  --go_opt=paths=source_relative \
  --go-grpc_out=services/link-rpc/pb \
  --go-grpc_opt=paths=source_relative \
  --zrpc_out=services/link-rpc \
  --proto_path=services/link-rpc/proto
```

Accept goctl's generated package, directory, and client names as the source of truth.
