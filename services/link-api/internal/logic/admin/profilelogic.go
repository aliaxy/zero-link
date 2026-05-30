// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-api/internal/apierror"
	"github.com/aliaxy/zero-link/services/link-api/internal/convert"
	"github.com/aliaxy/zero-link/services/link-api/internal/middleware"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/client/adminservice"

	"github.com/zeromicro/go-zero/core/logx"
)

// ProfileLogic coordinates authenticated administrator profile retrieval.
type ProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewProfileLogic creates administrator profile logic.
func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Profile returns the authenticated administrator profile.
func (l *ProfileLogic) Profile() (resp *types.ProfileResponse, err error) {
	subject, ok := middleware.AdminSubjectFromContext(l.ctx)
	if !ok {
		return nil, apierror.ErrUnauthenticated
	}

	rpcResp, err := l.svcCtx.AdminRPC.GetAdminProfile(l.ctx, &adminservice.GetAdminProfileRequest{
		AdminId: subject.ID,
	})
	if err != nil {
		return nil, apierror.FromRPCError(err)
	}

	return &types.ProfileResponse{
		Code:    "OK",
		Message: "ok",
		Data:    convert.AdminInfoFromRPC(rpcResp.Admin),
	}, nil
}
