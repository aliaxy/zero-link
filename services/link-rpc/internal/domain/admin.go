package domain

import (
	"context"
	"errors"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
)

const (
	// AdminStatusActive allows an administrator to use management APIs.
	AdminStatusActive = 1
	// AdminStatusDisabled blocks management authentication.
	AdminStatusDisabled = 2
)

var (
	// ErrInvalidArgument reports invalid management input.
	ErrInvalidArgument = errors.New("invalid argument")
	// ErrUnauthenticated reports invalid administrator credentials.
	ErrUnauthenticated = errors.New("unauthenticated")
	// ErrNotFound reports a missing or hidden management resource.
	ErrNotFound = errors.New("not found")
	// ErrConflict reports a duplicate management resource.
	ErrConflict = errors.New("conflict")
	// ErrPermissionDenied reports a disabled short link that cannot redirect.
	ErrPermissionDenied = errors.New("permission denied")
	// ErrGone reports an expired short link that can no longer redirect.
	ErrGone = errors.New("gone")
)

// AdminFinder finds administrator records by username.
type AdminFinder interface {
	FindOneByUsername(ctx context.Context, username string) (*model.AdminUser, error)
}

// AuthenticateAdmin validates administrator credentials and status.
func AuthenticateAdmin(ctx context.Context, finder AdminFinder, username, password string) (*model.AdminUser, error) {
	if username == "" || password == "" {
		return nil, ErrInvalidArgument
	}

	admin, err := finder.FindOneByUsername(ctx, username)
	if err != nil {
		return nil, ErrUnauthenticated
	}
	if admin.Status != AdminStatusActive {
		return nil, ErrUnauthenticated
	}
	if !VerifyPassword(admin.PasswordHash, password) {
		return nil, ErrUnauthenticated
	}

	return admin, nil
}

// GetActiveAdminProfile returns an active administrator profile.
func GetActiveAdminProfile(admin *model.AdminUser) (*model.AdminUser, error) {
	if admin == nil || admin.Status != AdminStatusActive {
		return nil, ErrNotFound
	}
	return admin, nil
}
