// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-api/internal/apierror"
	"github.com/aliaxy/zero-link/services/link-api/internal/convert"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/linkservice"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetShortLinkLogic coordinates management short-link detail retrieval.
type GetShortLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetShortLinkLogic creates short-link detail logic.
func NewGetShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShortLinkLogic {
	return &GetShortLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetShortLink returns one short link through link-rpc.
func (l *GetShortLinkLogic) GetShortLink(req *types.LinkIdRequest) (resp *types.ShortLinkResponse, err error) {
	rpcResp, err := l.svcCtx.LinkRPC.GetShortLink(l.ctx, &linkservice.GetShortLinkRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, apierror.FromRPCError(err)
	}

	return convert.OkShortLinkResponse(rpcResp.Link), nil
}
