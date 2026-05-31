// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

// Package redirect contains link-api short-link redirect logic.
package redirect

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-api/internal/apierror"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/client/linkservice"

	"github.com/zeromicro/go-zero/core/logx"
)

// RedirectLogic handles short-link redirect requests.
//
//nolint:revive
type RedirectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRedirectLogic creates a RedirectLogic.
func NewRedirectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedirectLogic {
	return &RedirectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Redirect resolves a short code and returns the origin URL for redirect.
func (l *RedirectLogic) Redirect(req *types.RedirectRequest) (*types.RedirectResponse, error) {
	resp, err := l.svcCtx.LinkRPC.ResolveShortLink(l.ctx, &linkservice.ResolveShortLinkRequest{
		Code: req.Code,
	})
	if err != nil {
		l.Errorw("redirect rpc failed",
			logx.Field("code", req.Code),
			logx.Field("error", err.Error()),
		)
		return nil, apierror.FromRPCError(err)
	}
	l.Infow("redirect resolved",
		logx.Field("code", req.Code),
		logx.Field("url", resp.OriginUrl),
	)
	return &types.RedirectResponse{OriginUrl: resp.OriginUrl}, nil
}
