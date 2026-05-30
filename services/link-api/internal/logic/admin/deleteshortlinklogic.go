// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-api/internal/apierror"
	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/client/linkservice"

	"github.com/zeromicro/go-zero/core/logx"
)

// DeleteShortLinkLogic coordinates management short-link deletion.
type DeleteShortLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteShortLinkLogic creates short-link deletion logic.
func NewDeleteShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteShortLinkLogic {
	return &DeleteShortLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteShortLink soft deletes a short link through link-rpc.
func (l *DeleteShortLinkLogic) DeleteShortLink(
	req *types.LinkIdRequest,
) (resp *types.DeleteShortLinkResponse, err error) {
	rpcResp, err := l.svcCtx.LinkRPC.DeleteShortLink(l.ctx, &linkservice.DeleteShortLinkRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, apierror.FromRPCError(err)
	}

	return &types.DeleteShortLinkResponse{
		Code:    "OK",
		Message: "ok",
		Data: types.DeleteShortLinkData{
			Id:      rpcResp.Id,
			Deleted: rpcResp.Deleted,
		},
	}, nil
}
