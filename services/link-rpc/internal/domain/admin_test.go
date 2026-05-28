package domain

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
)

type fakeAdminFinder struct {
	admin *model.AdminUser
	err   error
}

func (f fakeAdminFinder) FindOneByUsername(_ context.Context, _ string) (*model.AdminUser, error) {
	return f.admin, f.err
}

func TestAuthenticateAdmin(t *testing.T) {
	const seededHash = "$2a$10$UpD2JjqWVgQOatvqxd5H3OSQwzxC5o5gYf31R73AJIz.dQOAuKkBS"

	admin := &model.AdminUser{
		Id:           42,
		Username:     "admin",
		PasswordHash: seededHash,
		Status:       AdminStatusActive,
		CreatedAt:    time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC),
	}

	got, err := AuthenticateAdmin(context.Background(), fakeAdminFinder{admin: admin}, "admin", "zerolink")
	if err != nil {
		t.Fatalf("AuthenticateAdmin() error = %v", err)
	}
	if got.Id != 42 {
		t.Fatalf("admin ID = %d, want 42", got.Id)
	}
	if got.Username != "admin" {
		t.Fatalf("admin username = %q, want admin", got.Username)
	}
}

func TestAuthenticateAdminRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
	}{
		{name: "missing username", password: "zerolink"},
		{name: "missing password", username: "admin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AuthenticateAdmin(context.Background(), fakeAdminFinder{}, tt.username, tt.password)
			if !errors.Is(err, ErrInvalidArgument) {
				t.Fatalf("AuthenticateAdmin() error = %v, want ErrInvalidArgument", err)
			}
		})
	}
}

func TestAuthenticateAdminRejectsInvalidCredentials(t *testing.T) {
	const seededHash = "$2a$10$UpD2JjqWVgQOatvqxd5H3OSQwzxC5o5gYf31R73AJIz.dQOAuKkBS"

	admin := &model.AdminUser{
		Id:           42,
		Username:     "admin",
		PasswordHash: seededHash,
		Status:       AdminStatusActive,
	}

	_, err := AuthenticateAdmin(context.Background(), fakeAdminFinder{admin: admin}, "admin", "wrong-password")
	if !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("AuthenticateAdmin() error = %v, want ErrUnauthenticated", err)
	}
}

func TestAuthenticateAdminRejectsInactiveAdmin(t *testing.T) {
	const seededHash = "$2a$10$UpD2JjqWVgQOatvqxd5H3OSQwzxC5o5gYf31R73AJIz.dQOAuKkBS"

	admin := &model.AdminUser{
		Id:           42,
		Username:     "admin",
		PasswordHash: seededHash,
		Status:       AdminStatusDisabled,
	}

	_, err := AuthenticateAdmin(context.Background(), fakeAdminFinder{admin: admin}, "admin", "zerolink")
	if !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("AuthenticateAdmin() error = %v, want ErrUnauthenticated", err)
	}
}

func TestGetActiveAdminProfile(t *testing.T) {
	admin := &model.AdminUser{
		Id:        42,
		Username:  "admin",
		Status:    AdminStatusActive,
		CreatedAt: time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC),
	}

	got, err := GetActiveAdminProfile(admin)
	if err != nil {
		t.Fatalf("GetActiveAdminProfile() error = %v", err)
	}
	if got.Id != 42 {
		t.Fatalf("admin ID = %d, want 42", got.Id)
	}
}

func TestGetActiveAdminProfileRejectsInactiveAdmin(t *testing.T) {
	admin := &model.AdminUser{
		Id:       42,
		Username: "admin",
		Status:   AdminStatusDisabled,
	}

	_, err := GetActiveAdminProfile(admin)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetActiveAdminProfile() error = %v, want ErrNotFound", err)
	}
}
