// Package config defines link-rpc runtime configuration.
package config

import "github.com/zeromicro/go-zero/zrpc"

// MySQLConfig contains the MySQL dependency settings used by readiness checks.
type MySQLConfig struct {
	Endpoint   string
	Database   string
	User       string
	DataSource string
}

// RedisConfig contains the Redis dependency settings used by readiness checks.
type RedisConfig struct {
	Endpoint string
}

// DependenciesConfig contains external dependency settings for link-rpc.
type DependenciesConfig struct {
	MySQL MySQLConfig
	Redis RedisConfig
}

// Config contains the RPC server and dependency settings.
type Config struct {
	zrpc.RpcServerConf
	Dependencies DependenciesConfig
}
