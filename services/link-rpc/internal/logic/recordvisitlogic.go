package logic

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	"github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecordVisitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecordVisitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordVisitLogic {
	return &RecordVisitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RecordVisitLogic) RecordVisit(in *linkv1.RecordVisitRequest) (*linkv1.RecordVisitResponse, error) {
	// todo: add your logic here and delete this line

	return &linkv1.RecordVisitResponse{}, nil
}
