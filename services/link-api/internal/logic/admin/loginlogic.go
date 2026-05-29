// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

// Package admin contains link-api management request logic.
package admin

import (
	"context"
	"time"

	"github.com/aliaxy/zero-link/services/link-api/internal/apierror"
	"github.com/aliaxy/zero-link/services/link-api/internal/convert"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/linkservice"

	"github.com/zeromicro/go-zero/core/logx"
)

// LoginLogic coordinates administrator login.
type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewLoginLogic creates administrator login logic.
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login authenticates an administrator and returns a management token.
func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	rpcResp, err := l.svcCtx.LinkRPC.AuthenticateAdmin(l.ctx, &linkservice.AuthenticateAdminRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, apierror.FromRPCError(err)
	}

	token, expiresAt, err := l.svcCtx.TokenManager.Create(convert.AdminSubjectFromRPC(rpcResp.Admin))
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		Code:    "OK",
		Message: "ok",
		Data: types.LoginData{
			Token:     token,
			ExpiresAt: expiresAt.Format(time.RFC3339),
			Admin:     convert.AdminInfoFromRPC(rpcResp.Admin),
		},
	}, nil
}
