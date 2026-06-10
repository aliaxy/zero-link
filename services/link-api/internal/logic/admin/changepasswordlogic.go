// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-api/internal/apierror"
	"github.com/aliaxy/zero-link/services/link-api/internal/middleware"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/client/adminservice"

	"github.com/zeromicro/go-zero/core/logx"
)

// ChangePasswordLogic coordinates administrator password change and session revocation.
type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewChangePasswordLogic creates administrator password change logic.
func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChangePassword updates the administrator password and revokes all refresh tokens.
func (l *ChangePasswordLogic) ChangePassword(
	req *types.ChangePasswordRequest,
) (resp *types.ChangePasswordResponse, err error) {
	subject, ok := middleware.AdminSubjectFromContext(l.ctx)
	if !ok {
		return nil, apierror.ErrUnauthenticated
	}

	_, err = l.svcCtx.AdminRPC.ChangePassword(l.ctx, &adminservice.ChangePasswordRequest{
		AdminId:     subject.ID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		l.Infow("change password failed",
			logx.Field("admin_id", subject.ID),
			logx.Field("error", err.Error()),
		)
		return nil, apierror.FromRPCError(err)
	}

	// Revoke all refresh tokens so existing sessions must re-login.
	if revokeErr := l.svcCtx.RefreshTokenStore.RevokeAll(l.ctx, subject.ID); revokeErr != nil {
		l.Errorw("revoke refresh tokens failed",
			logx.Field("admin_id", subject.ID),
			logx.Field("error", revokeErr.Error()),
		)
	}

	l.Infow("change password success", logx.Field("admin_id", subject.ID))
	return &types.ChangePasswordResponse{Code: "OK", Message: "ok"}, nil
}
