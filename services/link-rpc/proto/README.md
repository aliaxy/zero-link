# RPC Contracts

This directory contains protobuf files.

The current protobuf contract includes only a readiness RPC contract. Stage 3 database foundation exists, but business RPC methods are reserved for the next Stage 3 implementation pass.

Proto files live under versioned package directories, such as `link/v1`, and should use absolute `go_package` values with explicit Go package aliases. For `link.v1`, use the Go alias `linkv1`.

If goctl generates awkward aliases for protobuf imports, normalize the generated service code to import the protobuf package as `linkv1`.
