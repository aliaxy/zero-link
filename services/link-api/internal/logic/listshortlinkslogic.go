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

// ListShortLinksLogic coordinates management short-link listing.
type ListShortLinksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListShortLinksLogic creates short-link listing logic.
func NewListShortLinksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListShortLinksLogic {
	return &ListShortLinksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListShortLinks returns paginated short-link summaries through link-rpc.
func (l *ListShortLinksLogic) ListShortLinks(
	req *types.ListShortLinksRequest,
) (resp *types.ListShortLinksResponse, err error) {
	rpcResp, err := l.svcCtx.LinkRPC.ListShortLinks(l.ctx, &linkservice.ListShortLinksRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.Status,
		Keyword:  req.Keyword,
	})
	if err != nil {
		return nil, err
	}

	items := make([]types.ShortLinkSummary, 0, len(rpcResp.Items))
	for _, item := range rpcResp.Items {
		items = append(items, shortLinkSummaryFromRPC(item))
	}

	return &types.ListShortLinksResponse{
		Code:    "OK",
		Message: "ok",
		Data: types.ListShortLinksData{
			Items:    items,
			Page:     rpcResp.Page,
			PageSize: rpcResp.PageSize,
			Total:    rpcResp.Total,
		},
	}, nil
}
