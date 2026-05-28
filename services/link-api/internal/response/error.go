// Package response centralizes link-api response mapping.
package response

import (
	"context"
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorEnvelope is the stable management API error response body.
type ErrorEnvelope struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler maps service errors to management API response envelopes.
func ErrorHandler(ctx context.Context, err error) (int, any) {
	_ = ctx

	if st, ok := status.FromError(err); ok {
		return grpcError(st)
	}
	if errors.Is(err, ErrUnauthenticated) {
		return http.StatusUnauthorized, ErrorEnvelope{
			Code:    "UNAUTHENTICATED",
			Message: "missing or invalid bearer token",
		}
	}

	return http.StatusInternalServerError, ErrorEnvelope{
		Code:    "INTERNAL",
		Message: "internal error",
	}
}

// ErrUnauthenticated reports a missing or invalid management identity.
var ErrUnauthenticated = errors.New("unauthenticated")

func grpcError(st *status.Status) (int, ErrorEnvelope) {
	switch st.Code() {
	case codes.InvalidArgument:
		return http.StatusBadRequest, ErrorEnvelope{Code: "INVALID_ARGUMENT", Message: st.Message()}
	case codes.Unauthenticated:
		return http.StatusUnauthorized, ErrorEnvelope{Code: "UNAUTHENTICATED", Message: st.Message()}
	case codes.PermissionDenied:
		return http.StatusForbidden, ErrorEnvelope{Code: "PERMISSION_DENIED", Message: st.Message()}
	case codes.NotFound:
		return http.StatusNotFound, ErrorEnvelope{Code: "NOT_FOUND", Message: st.Message()}
	case codes.AlreadyExists:
		return http.StatusConflict, ErrorEnvelope{Code: "CONFLICT", Message: st.Message()}
	default:
		return http.StatusInternalServerError, ErrorEnvelope{Code: "INTERNAL", Message: "internal error"}
	}
}
