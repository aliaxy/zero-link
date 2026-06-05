// Package rpcerror converts domain errors to gRPC status errors.
package rpcerror

import (
	"errors"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToRPC maps a domain error to the appropriate gRPC status error.
func ToRPC(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, domain.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrConflict):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrPermissionDenied):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, domain.ErrGone):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
