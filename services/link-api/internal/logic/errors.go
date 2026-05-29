package logic

import "github.com/aliaxy/zero-link/services/link-api/internal/response"

func errUnauthenticated() error {
	return response.ErrUnauthenticated
}

// fromRPCError passes an RPC error through to the framework's ErrorHandler.
// link-rpc wraps all domain errors in gRPC status codes via rpcError();
// response.ErrorHandler translates those codes to HTTP envelopes.
func fromRPCError(err error) error {
	return err
}
