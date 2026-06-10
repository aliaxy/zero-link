package adminservicelogic

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/rpcerror"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

// ChangePasswordLogic coordinates administrator password change.
type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewChangePasswordLogic creates administrator password change logic.
func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChangePassword verifies the old password and stores the new bcrypt hash.
func (l *ChangePasswordLogic) ChangePassword(in *linkv1.ChangePasswordRequest) (*linkv1.ChangePasswordResponse, error) {
	err := domain.ChangePassword(
		l.ctx,
		l.svcCtx.AdminUserModel,
		in.GetAdminId(),
		in.GetOldPassword(),
		in.GetNewPassword(),
	)
	if err != nil {
		l.Infow("change password failed", logx.Field("admin_id", in.GetAdminId()))
		return nil, rpcerror.ToRPC(err)
	}

	l.Infow("change password success", logx.Field("admin_id", in.GetAdminId()))
	return &linkv1.ChangePasswordResponse{}, nil
}
