// Package config defines link-rpc runtime configuration.
package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

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

// AnalyticsConfig contains analytics-specific settings.
type AnalyticsConfig struct {
	IPSalt string
}

// RetentionConfig controls data lifecycle cleanup behaviour.
type RetentionConfig struct {
	VisitEventRetentionDays int
	ShortLinkRetentionDays  int
	DailyStatRetentionDays  int
	CleanupBatchSize        int
}

// CuckooConfig controls the in-process cuckoo filter for cache penetration defence.
type CuckooConfig struct {
	Capacity uint
}

// Config contains the RPC server and dependency settings.
type Config struct {
	zrpc.RpcServerConf
	Dependencies DependenciesConfig
	CacheRedis   cache.CacheConf
	Analytics    AnalyticsConfig
	Retention    RetentionConfig
	Cuckoo       CuckooConfig
}
