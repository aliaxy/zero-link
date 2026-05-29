// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package redirect

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/linkservice"

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
		return nil, fromRPCError(err)
	}
	return &types.RedirectResponse{OriginUrl: resp.OriginUrl}, nil
}
