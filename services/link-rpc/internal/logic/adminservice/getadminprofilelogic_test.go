package adminservicelogic

import (
	"context"
	"testing"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeProfileAdminModel struct {
	model.AdminUserModel
	admin *model.AdminUser
	err   error
}

func (f fakeProfileAdminModel) FindOne(_ context.Context, _ int64) (*model.AdminUser, error) {
	return f.admin, f.err
}

func newProfileLogic(m model.AdminUserModel) *GetAdminProfileLogic {
	return NewGetAdminProfileLogic(context.Background(), &svc.ServiceContext{
		AdminUserModel: m,
	})
}

func TestGetAdminProfile_Found(t *testing.T) {
	t.Parallel()
	admin := activeAdmin()
	resp, err := newProfileLogic(fakeProfileAdminModel{admin: admin}).
		GetAdminProfile(&linkv1.GetAdminProfileRequest{AdminId: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Admin.Id != 42 {
		t.Fatalf("want id 42, got %d", resp.Admin.Id)
	}
	if resp.Admin.Username != "admin" {
		t.Fatalf("want username admin, got %q", resp.Admin.Username)
	}
}

func TestGetAdminProfile_NotFound(t *testing.T) {
	t.Parallel()
	m := fakeProfileAdminModel{err: model.ErrNotFound}
	_, err := newProfileLogic(m).GetAdminProfile(&linkv1.GetAdminProfileRequest{AdminId: 99})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("want NotFound, got %v", err)
	}
}

func TestGetAdminProfile_DisabledAdminHidden(t *testing.T) {
	t.Parallel()
	admin := activeAdmin()
	admin.Status = domain.AdminStatusDisabled
	m := fakeProfileAdminModel{admin: admin}
	_, err := newProfileLogic(m).GetAdminProfile(&linkv1.GetAdminProfileRequest{AdminId: 42})
	// GetActiveAdminProfile returns ErrNotFound for disabled admins
	if status.Code(err) != codes.NotFound {
		t.Fatalf("want NotFound for disabled admin, got %v", err)
	}
}
