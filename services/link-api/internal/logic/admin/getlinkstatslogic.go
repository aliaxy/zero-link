// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/linkservice"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetLinkStatsLogic handles daily stat query for a short link.
type GetLinkStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetLinkStatsLogic creates a GetLinkStatsLogic.
func NewGetLinkStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLinkStatsLogic {
	return &GetLinkStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetLinkStats returns daily PV/UV stats for a short link within a date range.
func (l *GetLinkStatsLogic) GetLinkStats(req *types.GetLinkStatsRequest) (resp *types.GetLinkStatsResponse, err error) {
	rpcResp, err := l.svcCtx.LinkRPC.GetLinkStats(l.ctx, &linkservice.GetLinkStatsRequest{
		LinkId: req.Id,
		From:   req.From,
		To:     req.To,
	})
	if err != nil {
		return nil, fromRPCError(err)
	}

	items := make([]types.DailyStatItem, 0, len(rpcResp.Items))
	for _, s := range rpcResp.Items {
		items = append(items, types.DailyStatItem{
			StatDate: s.StatDate,
			Pv:       s.Pv,
			Uv:       s.Uv,
		})
	}

	return &types.GetLinkStatsResponse{
		Code:    "SUCCESS",
		Message: "ok",
		Data: types.GetLinkStatsData{
			LinkId: rpcResp.LinkId,
			Items:  items,
		},
	}, nil
}
