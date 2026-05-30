// Package healthservicelogic contains link-rpc health service request logic.
package healthservicelogic

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

const dependencyDialTimeout = 2 * time.Second

// CheckLogic handles RPC readiness checks.
type CheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewCheckLogic creates a readiness check logic instance.
func NewCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckLogic {
	return &CheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Check reports whether link-rpc dependencies are reachable.
func (l *CheckLogic) Check(_ *linkv1.CheckRequest) (*linkv1.CheckResponse, error) {
	if err := validateEndpoint("mysql", l.svcCtx.Config.Dependencies.MySQL.Endpoint); err != nil {
		return &linkv1.CheckResponse{
			Ok:      false,
			Message: err.Error(),
		}, nil
	}

	if err := validateEndpoint("redis", l.svcCtx.Config.Dependencies.Redis.Endpoint); err != nil {
		return &linkv1.CheckResponse{
			Ok:      false,
			Message: err.Error(),
		}, nil
	}

	if err := dialEndpoint(l.ctx, "mysql", l.svcCtx.Config.Dependencies.MySQL.Endpoint); err != nil {
		return &linkv1.CheckResponse{
			Ok:      false,
			Message: err.Error(),
		}, nil
	}

	if err := dialEndpoint(l.ctx, "redis", l.svcCtx.Config.Dependencies.Redis.Endpoint); err != nil {
		return &linkv1.CheckResponse{
			Ok:      false,
			Message: err.Error(),
		}, nil
	}

	return &linkv1.CheckResponse{
		Ok:      true,
		Message: "rpc dependencies ready",
	}, nil
}

func validateEndpoint(name, endpoint string) error {
	if endpoint == "" {
		return fmt.Errorf("%s endpoint is empty", name)
	}

	return nil
}

func dialEndpoint(ctx context.Context, name, endpoint string) error {
	dialer := net.Dialer{Timeout: dependencyDialTimeout}
	conn, err := dialer.DialContext(ctx, "tcp", endpoint)
	if err != nil {
		return fmt.Errorf("%s unavailable: %w", name, err)
	}
	defer func() {
		_ = conn.Close()
	}()

	return nil
}
