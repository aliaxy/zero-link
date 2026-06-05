package analyticsservicelogic

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/domain"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/rpcerror"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetLinkStatsLogic handles daily stat query RPC.
type GetLinkStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewGetLinkStatsLogic creates a GetLinkStatsLogic.
func NewGetLinkStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLinkStatsLogic {
	return &GetLinkStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetLinkStats returns daily PV/UV stats for a link within a date range.
func (l *GetLinkStatsLogic) GetLinkStats(in *linkv1.GetLinkStatsRequest) (*linkv1.GetLinkStatsResponse, error) {
	from, to, err := parseDateRange(in.From, in.To)
	if err != nil {
		return nil, rpcerror.ToRPC(domain.ErrInvalidArgument)
	}

	stats, err := l.svcCtx.DailyStatModel.FindByLinkIDAndDateRange(l.ctx, in.LinkId, from, to)
	if err != nil {
		return nil, rpcerror.ToRPC(err)
	}

	items := make([]*linkv1.DailyStat, 0, len(stats))
	for _, s := range stats {
		items = append(items, &linkv1.DailyStat{
			StatDate: s.StatDate.Format(dateFormat),
			Pv:       s.Pv,
			Uv:       s.Uv,
		})
	}

	return &linkv1.GetLinkStatsResponse{
		LinkId: in.LinkId,
		Items:  items,
	}, nil
}
