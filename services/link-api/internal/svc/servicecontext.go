// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

// Package svc wires link-api service dependencies.
package svc

import (
	"github.com/aliaxy/zero-link/services/link-api/internal/auth"
	"github.com/aliaxy/zero-link/services/link-api/internal/config"
	"github.com/aliaxy/zero-link/services/link-api/internal/middleware"
	"github.com/aliaxy/zero-link/services/link-rpc/client/adminservice"
	"github.com/aliaxy/zero-link/services/link-rpc/client/analyticsservice"
	"github.com/aliaxy/zero-link/services/link-rpc/client/healthservice"
	"github.com/aliaxy/zero-link/services/link-rpc/client/linkservice"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// ServiceContext holds dependencies shared by link-api handlers.
type ServiceContext struct {
	Config                config.Config
	AuthMiddleware        rest.Middleware
	AnalyticsMiddleware   rest.Middleware
	IPRateLimitMiddleware    rest.Middleware
	LoginRateLimitMiddleware rest.Middleware
	TokenManager          *auth.TokenManager
	Redis                 *redis.Redis
	HealthRPC             healthservice.HealthService
	AdminRPC              adminservice.AdminService
	LinkRPC               linkservice.LinkService
	AnalyticsRPC          analyticsservice.AnalyticsService
}

// NewServiceContext creates a link-api service context.
func NewServiceContext(c config.Config) *ServiceContext {
	tokenManager := auth.NewTokenManager(auth.Config{
		Secret:          c.Auth.Secret,
		TokenTTLSeconds: c.Auth.TokenTTLSeconds,
	})

	rpcClient := zrpc.MustNewClient(c.LinkRPC)
	analyticsRPC := analyticsservice.NewAnalyticsService(rpcClient)
	redisClient := redis.MustNewRedis(c.Redis)

	return &ServiceContext{
		Config:              c,
		TokenManager:        tokenManager,
		Redis:               redisClient,
		AuthMiddleware:      middleware.NewAuthMiddleware(tokenManager).Handle,
		AnalyticsMiddleware: middleware.NewAnalyticsMiddleware(analyticsRPC).Handle,
		IPRateLimitMiddleware: middleware.NewIPRateLimitMiddleware(
			redisClient, 1, c.RateLimit.RedirectPerIPPerSecond, "rl:redirect:ip:",
		).Handle,
		LoginRateLimitMiddleware: middleware.NewLoginRateLimitMiddleware(
			redisClient, c.RateLimit.LoginPerIPPerMinute,
		).Handle,
		HealthRPC:    healthservice.NewHealthService(rpcClient),
		AdminRPC:     adminservice.NewAdminService(rpcClient),
		LinkRPC:      linkservice.NewLinkService(rpcClient),
		AnalyticsRPC: analyticsRPC,
	}
}
