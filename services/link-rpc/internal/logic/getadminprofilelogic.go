package logic

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	"github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetAdminProfileLogic coordinates administrator profile retrieval.
type GetAdminProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewGetAdminProfileLogic creates administrator profile logic.
func NewGetAdminProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminProfileLogic {
	return &GetAdminProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetAdminProfile returns an active administrator profile.
func (l *GetAdminProfileLogic) GetAdminProfile(
	in *linkv1.GetAdminProfileRequest,
) (*linkv1.GetAdminProfileResponse, error) {
	admin, err := l.svcCtx.AdminUserModel.FindOne(l.ctx, in.GetAdminId())
	if err != nil {
		return nil, rpcError(domain.ErrNotFound)
	}
	admin, err = domain.GetActiveAdminProfile(admin)
	if err != nil {
		return nil, rpcError(err)
	}

	return &linkv1.GetAdminProfileResponse{
		Admin: adminProfile(admin),
	}, nil
}
