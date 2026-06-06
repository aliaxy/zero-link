// Package apierror centralizes link-api error mapping.
package apierror

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

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
	if isParseError(err) {
		return http.StatusBadRequest, ErrorEnvelope{
			Code:    "BAD_REQUEST",
			Message: "invalid request parameters",
		}
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

// RedirectError writes a plain HTTP error for redirect path failures.
// The redirect path does not use JSON envelopes.
func RedirectError(w http.ResponseWriter, err error) {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.NotFound:
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		case codes.PermissionDenied:
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		case codes.FailedPrecondition:
			http.Error(w, http.StatusText(http.StatusGone), http.StatusGone)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// FromRPCError passes an RPC error through to the framework's ErrorHandler.
// link-rpc wraps all domain errors in gRPC status codes via rpcError();
// apierror.ErrorHandler translates those codes to HTTP envelopes.
func FromRPCError(err error) error {
	return err
}

// isParseError reports whether err originates from go-zero's httpx.Parse:
// missing required fields ("is not set") or malformed JSON bodies.
func isParseError(err error) bool {
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, " is not set") ||
		strings.Contains(msg, "type mismatch") ||
		strings.Contains(msg, "out of range")
}

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
	case codes.FailedPrecondition:
		return http.StatusGone, ErrorEnvelope{Code: "GONE", Message: st.Message()}
	case codes.Unavailable:
		return http.StatusServiceUnavailable, ErrorEnvelope{
			Code:    "SERVICE_UNAVAILABLE",
			Message: "service temporarily unavailable",
		}
	default:
		return http.StatusInternalServerError, ErrorEnvelope{Code: "INTERNAL", Message: "internal error"}
	}
}
