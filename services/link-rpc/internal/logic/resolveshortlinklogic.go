package logic

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	"github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResolveShortLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResolveShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResolveShortLinkLogic {
	return &ResolveShortLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResolveShortLinkLogic) ResolveShortLink(in *linkv1.ResolveShortLinkRequest) (*linkv1.ResolveShortLinkResponse, error) {
	// todo: add your logic here and delete this line

	return &linkv1.ResolveShortLinkResponse{}, nil
}
