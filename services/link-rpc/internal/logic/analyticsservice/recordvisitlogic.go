package analyticsservicelogic

import (
	"context"
	"errors"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/metrics"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/rpcerror"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

// RecordVisitLogic handles visit event recording RPC.
type RecordVisitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewRecordVisitLogic creates a RecordVisitLogic.
func NewRecordVisitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordVisitLogic {
	return &RecordVisitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RecordVisit records a visit event and upserts the daily stat.
func (l *RecordVisitLogic) RecordVisit(in *linkv1.RecordVisitRequest) (*linkv1.RecordVisitResponse, error) {
	visitedAt, err := time.Parse(time.RFC3339Nano, in.VisitedAt)
	if err != nil {
		visitedAt = time.Now().UTC()
	}

	link, err := l.svcCtx.ShortLinkModel.FindOneByCode(l.ctx, in.Code)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &linkv1.RecordVisitResponse{}, nil
		}
		return nil, rpcerror.ToRPC(err)
	}

	ipHash := hashIP(in.Ip, l.svcCtx.Config.Analytics.IPSalt)
	statDate := visitedAt.UTC().Format(dateFormat)

	alreadyVisited, err := l.svcCtx.VisitEventModel.HasVisitedToday(l.ctx, link.Id, ipHash, statDate)
	if err != nil {
		l.Errorw("check uv failed, defaulting to no-uv",
			logx.Field("link_id", link.Id),
			logx.Field("error", err.Error()),
		)
		alreadyVisited = true
	}

	event := &model.VisitEvent{
		LinkId:    link.Id,
		Code:      in.Code,
		VisitedAt: visitedAt,
		IpHash:    ipHash,
		UserAgent: in.UserAgent,
		Referer:   in.Referer,
		Device:    detectDevice(in.UserAgent),
	}
	if _, err := l.svcCtx.VisitEventModel.Insert(l.ctx, event); err != nil {
		return nil, rpcerror.ToRPC(err)
	}

	if err := l.svcCtx.DailyStatModel.UpsertStats(l.ctx, link.Id, statDate, !alreadyVisited); err != nil {
		l.Errorw("upsert daily stat failed",
			logx.Field("link_id", link.Id),
			logx.Field("code", in.Code),
			logx.Field("stat_date", statDate),
			logx.Field("error", err.Error()),
		)
		metrics.AnalyticsEventsTotal.Inc("error")
		return &linkv1.RecordVisitResponse{}, nil
	}

	metrics.AnalyticsEventsTotal.Inc("success")
	return &linkv1.RecordVisitResponse{}, nil
}
