package logic

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

// ResolveShortLinkLogic handles short-link resolution RPC.
type ResolveShortLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewResolveShortLinkLogic creates a ResolveShortLinkLogic.
func NewResolveShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResolveShortLinkLogic {
	return &ResolveShortLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ResolveShortLink resolves a short code to its origin URL.
func (l *ResolveShortLinkLogic) ResolveShortLink(
	_ *linkv1.ResolveShortLinkRequest,
) (*linkv1.ResolveShortLinkResponse, error) {
	// todo: add your logic here and delete this line

	return &linkv1.ResolveShortLinkResponse{}, nil
}
