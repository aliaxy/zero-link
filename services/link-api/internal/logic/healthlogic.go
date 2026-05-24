// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

// Package logic contains link-api request logic.
package logic

import (
	"context"

	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// HealthLogic handles API liveness checks.
type HealthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewHealthLogic creates a liveness check logic instance.
func NewHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthLogic {
	return &HealthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Health reports whether the API process is alive.
func (l *HealthLogic) Health() (resp *types.HealthResponse, err error) {
	return &types.HealthResponse{
		Status:  "ok",
		Message: "api alive",
	}, nil
}
