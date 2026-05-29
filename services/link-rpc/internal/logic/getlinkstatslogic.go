package logic

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	"github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLinkStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLinkStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLinkStatsLogic {
	return &GetLinkStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLinkStatsLogic) GetLinkStats(in *linkv1.GetLinkStatsRequest) (*linkv1.GetLinkStatsResponse, error) {
	// todo: add your logic here and delete this line

	return &linkv1.GetLinkStatsResponse{}, nil
}
