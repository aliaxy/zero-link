// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

// Package config defines link-api runtime configuration.
package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// AuthConfig contains management JWT settings.
type AuthConfig struct {
	Secret          string
	TokenTTLSeconds int64
}

// Config contains the HTTP server and upstream RPC client settings.
type Config struct {
	rest.RestConf
	Auth    AuthConfig
	LinkRPC zrpc.RpcClientConf
}
