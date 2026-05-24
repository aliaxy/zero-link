// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"

	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/linkclient"

	"github.com/zeromicro/go-zero/core/logx"
)

// ReadyLogic handles API readiness checks.
type ReadyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewReadyLogic creates a readiness check logic instance.
func NewReadyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadyLogic {
	return &ReadyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Ready reports whether the API and its RPC dependency are ready.
func (l *ReadyLogic) Ready() (resp *types.ReadyResponse, err error) {
	check, err := l.svcCtx.LinkRPC.Check(l.ctx, &linkclient.CheckRequest{})
	if err != nil {
		return nil, fmt.Errorf("check link rpc readiness: %w", err)
	}
	if !check.Ok {
		return nil, fmt.Errorf("link rpc not ready: %s", check.Message)
	}

	return &types.ReadyResponse{
		Status:  "ok",
		Message: "api and rpc ready",
	}, nil
}
