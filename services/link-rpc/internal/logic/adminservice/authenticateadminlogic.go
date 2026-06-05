// Package adminservicelogic contains link-rpc admin service request logic.
package adminservicelogic

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/rpcerror"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

// AuthenticateAdminLogic coordinates administrator authentication.
type AuthenticateAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewAuthenticateAdminLogic creates administrator authentication logic.
func NewAuthenticateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthenticateAdminLogic {
	return &AuthenticateAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AuthenticateAdmin validates administrator credentials.
func (l *AuthenticateAdminLogic) AuthenticateAdmin(
	in *linkv1.AuthenticateAdminRequest,
) (*linkv1.AuthenticateAdminResponse, error) {
	admin, err := domain.AuthenticateAdmin(l.ctx, l.svcCtx.AdminUserModel, in.GetUsername(), in.GetPassword())
	if err != nil {
		l.Infow("authenticate failed", logx.Field("username", in.GetUsername()))
		return nil, rpcerror.ToRPC(err)
	}

	l.Infow("authenticate success",
		logx.Field("admin_id", admin.Id),
		logx.Field("username", admin.Username),
	)
	return &linkv1.AuthenticateAdminResponse{
		Admin: adminProfile(admin),
	}, nil
}

func adminProfile(admin *model.AdminUser) *linkv1.AdminProfile {
	return &linkv1.AdminProfile{
		Id:        admin.Id,
		Username:  admin.Username,
		Status:    admin.Status,
		CreatedAt: admin.CreatedAt.Format(timeFormat),
	}
}
