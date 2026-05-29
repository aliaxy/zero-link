package apierror

import (
	"context"
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorHandlerMapsGRPCErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{
			name:       "invalid argument",
			err:        status.Error(codes.InvalidArgument, "bad input"),
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_ARGUMENT",
		},
		{
			name:       "unauthenticated",
			err:        status.Error(codes.Unauthenticated, "invalid credentials"),
			wantStatus: http.StatusUnauthorized,
			wantCode:   "UNAUTHENTICATED",
		},
		{
			name:       "not found",
			err:        status.Error(codes.NotFound, "not found"),
			wantStatus: http.StatusNotFound,
			wantCode:   "NOT_FOUND",
		},
		{
			name:       "conflict",
			err:        status.Error(codes.AlreadyExists, "conflict"),
			wantStatus: http.StatusConflict,
			wantCode:   "CONFLICT",
		},
		{
			name:       "permission denied",
			err:        status.Error(codes.PermissionDenied, "disabled"),
			wantStatus: http.StatusForbidden,
			wantCode:   "PERMISSION_DENIED",
		},
		{
			name:       "gone",
			err:        status.Error(codes.FailedPrecondition, "expired"),
			wantStatus: http.StatusGone,
			wantCode:   "GONE",
		},
		{
			name:       "internal",
			err:        status.Error(codes.Internal, "db failed"),
			wantStatus: http.StatusInternalServerError,
			wantCode:   "INTERNAL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, gotEnvelope := ErrorHandler(context.Background(), tt.err)
			if gotStatus != tt.wantStatus {
				t.Fatalf("status = %d, want %d", gotStatus, tt.wantStatus)
			}

			envelope, ok := gotEnvelope.(ErrorEnvelope)
			if !ok {
				t.Fatalf("envelope type = %T, want ErrorEnvelope", gotEnvelope)
			}
			if envelope.Code != tt.wantCode {
				t.Fatalf("code = %q, want %q", envelope.Code, tt.wantCode)
			}
		})
	}
}
