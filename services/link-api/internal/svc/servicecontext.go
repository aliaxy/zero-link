// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

// Package svc wires link-api service dependencies.
package svc

import (
	"github.com/aliaxy/zero-link/services/link-api/internal/config"
	"github.com/aliaxy/zero-link/services/link-rpc/linkclient"

	"github.com/zeromicro/go-zero/zrpc"
)

// ServiceContext holds dependencies shared by link-api handlers.
type ServiceContext struct {
	Config  config.Config
	LinkRPC linkclient.Link
}

// NewServiceContext creates a link-api service context.
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		LinkRPC: linkclient.NewLink(zrpc.MustNewClient(c.LinkRPC)),
	}
}
