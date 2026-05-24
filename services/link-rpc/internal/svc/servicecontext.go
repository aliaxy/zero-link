// Package svc wires link-rpc service dependencies.
package svc

import "github.com/aliaxy/zero-link/services/link-rpc/internal/config"

// ServiceContext holds dependencies shared by link-rpc logic.
type ServiceContext struct {
	Config config.Config
}

// NewServiceContext creates a link-rpc service context.
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
