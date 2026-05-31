package adminservicelogic

import (
	"context"
	"testing"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// seededHash is the bcrypt hash of "zerolink" used across tests.
const seededHash = "$2a$10$UpD2JjqWVgQOatvqxd5H3OSQwzxC5o5gYf31R73AJIz.dQOAuKkBS"

type fakeAuthAdminModel struct {
	model.AdminUserModel
	admin *model.AdminUser
	err   error
}

func (f fakeAuthAdminModel) FindOneByUsername(_ context.Context, _ string) (*model.AdminUser, error) {
	return f.admin, f.err
}

func activeAdmin() *model.AdminUser {
	return &model.AdminUser{
		Id:           42,
		Username:     "admin",
		PasswordHash: seededHash,
		Status:       domain.AdminStatusActive,
		CreatedAt:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func newAuthLogic(m model.AdminUserModel) *AuthenticateAdminLogic {
	return NewAuthenticateAdminLogic(context.Background(), &svc.ServiceContext{
		AdminUserModel: m,
	})
}

func TestAuthenticateAdmin_ValidCredentials(t *testing.T) {
	resp, err := newAuthLogic(fakeAuthAdminModel{admin: activeAdmin()}).
		AuthenticateAdmin(&linkv1.AuthenticateAdminRequest{
			Username: "admin",
			Password: "zerolink",
		})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Admin.Id != 42 {
		t.Fatalf("want admin id 42, got %d", resp.Admin.Id)
	}
	if resp.Admin.Username != "admin" {
		t.Fatalf("want username admin, got %q", resp.Admin.Username)
	}
}

func TestAuthenticateAdmin_WrongPassword(t *testing.T) {
	_, err := newAuthLogic(fakeAuthAdminModel{admin: activeAdmin()}).
		AuthenticateAdmin(&linkv1.AuthenticateAdminRequest{
			Username: "admin",
			Password: "wrong",
		})
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("want Unauthenticated, got %v", err)
	}
}

func TestAuthenticateAdmin_AdminNotFound(t *testing.T) {
	_, err := newAuthLogic(fakeAuthAdminModel{err: model.ErrNotFound}).
		AuthenticateAdmin(&linkv1.AuthenticateAdminRequest{
			Username: "nobody",
			Password: "zerolink",
		})
	// domain.AuthenticateAdmin returns ErrUnauthenticated when admin is not found
	// to prevent username enumeration
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("want Unauthenticated for missing admin, got %v", err)
	}
}

func TestAuthenticateAdmin_DisabledAdmin(t *testing.T) {
	admin := activeAdmin()
	admin.Status = domain.AdminStatusDisabled
	_, err := newAuthLogic(fakeAuthAdminModel{admin: admin}).
		AuthenticateAdmin(&linkv1.AuthenticateAdminRequest{
			Username: "admin",
			Password: "zerolink",
		})
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("want Unauthenticated for disabled admin, got %v", err)
	}
}

func TestAuthenticateAdmin_MissingCredentials(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		username string
		password string
	}{
		{name: "missing username", username: "", password: "zerolink"},
		{name: "missing password", username: "admin", password: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := newAuthLogic(fakeAuthAdminModel{admin: activeAdmin()}).
				AuthenticateAdmin(&linkv1.AuthenticateAdminRequest{
					Username: tt.username,
					Password: tt.password,
				})
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("want InvalidArgument, got %v", err)
			}
		})
	}
}
