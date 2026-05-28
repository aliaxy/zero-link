// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

// Package svc wires link-api service dependencies.
package svc

import (
	"github.com/aliaxy/zero-link/services/link-api/internal/auth"
	"github.com/aliaxy/zero-link/services/link-api/internal/config"
	"github.com/aliaxy/zero-link/services/link-api/internal/middleware"
	"github.com/aliaxy/zero-link/services/link-rpc/linkservice"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// ServiceContext holds dependencies shared by link-api handlers.
type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware
	TokenManager   *auth.TokenManager
	LinkRPC        linkservice.LinkService
}

// NewServiceContext creates a link-api service context.
func NewServiceContext(c config.Config) *ServiceContext {
	tokenManager := auth.NewTokenManager(auth.Config{
		Secret:          c.Auth.Secret,
		TokenTTLSeconds: c.Auth.TokenTTLSeconds,
	})

	return &ServiceContext{
		Config:         c,
		TokenManager:   tokenManager,
		AuthMiddleware: middleware.NewAuthMiddleware(tokenManager).Handle,
		LinkRPC:        linkservice.NewLinkService(zrpc.MustNewClient(c.LinkRPC)),
	}
}
