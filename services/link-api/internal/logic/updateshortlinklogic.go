// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/linkservice"

	"github.com/zeromicro/go-zero/core/logx"
)

// UpdateShortLinkLogic coordinates management short-link updates.
type UpdateShortLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateShortLinkLogic creates short-link update logic.
func NewUpdateShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateShortLinkLogic {
	return &UpdateShortLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateShortLink updates mutable short-link fields through link-rpc.
func (l *UpdateShortLinkLogic) UpdateShortLink(
	req *types.UpdateShortLinkRequest,
) (resp *types.ShortLinkResponse, err error) {
	rpcResp, err := l.svcCtx.LinkRPC.UpdateShortLink(l.ctx, &linkservice.UpdateShortLinkRequest{
		Id:          req.Id,
		OriginUrl:   req.OriginUrl,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		ExpireAt:    req.ExpireAt,
	})
	if err != nil {
		return nil, fromRPCError(err)
	}

	return okShortLinkResponse(rpcResp.Link), nil
}
