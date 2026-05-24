// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

// Package config defines link-api runtime configuration.
package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config contains the HTTP server and upstream RPC client settings.
type Config struct {
	rest.RestConf
	LinkRPC zrpc.RpcClientConf
}
