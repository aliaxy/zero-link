// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"time"

	"github.com/aliaxy/zero-link/services/link-api/internal/apierror"
	"github.com/aliaxy/zero-link/services/link-api/internal/auth"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/client/adminservice"

	"github.com/zeromicro/go-zero/core/logx"
)

// RefreshTokenLogic rotates a refresh token and issues a new access token.
type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRefreshTokenLogic creates refresh token rotation logic.
func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// RefreshToken rotates the refresh token and returns a new access token.
func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenRequest) (resp *types.RefreshTokenResponse, err error) {
	newRefresh, adminID, err := l.svcCtx.RefreshTokenStore.Rotate(l.ctx, req.RefreshToken)
	if err != nil {
		l.Infow("refresh token invalid or expired")
		return nil, apierror.ErrUnauthenticated
	}

	profileResp, err := l.svcCtx.AdminRPC.GetAdminProfile(l.ctx, &adminservice.GetAdminProfileRequest{
		AdminId: adminID,
	})
	if err != nil {
		l.Errorw("get admin profile failed after token rotate",
			logx.Field("admin_id", adminID),
			logx.Field("error", err.Error()),
		)
		return nil, apierror.ErrUnauthenticated
	}

	accessToken, accessExpiresAt, err := l.svcCtx.TokenManager.Create(auth.AdminSubject{
		ID:       profileResp.Admin.Id,
		Username: profileResp.Admin.Username,
	})
	if err != nil {
		return nil, err
	}

	l.Infow("token refreshed", logx.Field("admin_id", adminID))

	return &types.RefreshTokenResponse{
		Code:    "OK",
		Message: "ok",
		Data: types.RefreshTokenData{
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  accessExpiresAt.Format(time.RFC3339),
			RefreshToken:          newRefresh,
			RefreshTokenExpiresAt: auth.ExpiresAt().Format(time.RFC3339),
		},
	}, nil
}
