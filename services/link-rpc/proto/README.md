# RPC Contracts

This directory contains protobuf files.

Stage 2 includes only a readiness RPC contract. Business RPC methods are reserved for later stages.

Proto files live under versioned package directories, such as `link/v1`, and should use absolute `go_package` values with explicit Go package aliases. For `link.v1`, use the Go alias `linkv1`.

If goctl generates awkward aliases for protobuf imports, normalize the generated service code to import the protobuf package as `linkv1`.
