// Package svc wires link-rpc service dependencies.
package svc

import (
	"github.com/aliaxy/zero-link/services/link-rpc/internal/config"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ServiceContext holds dependencies shared by link-rpc logic.
type ServiceContext struct {
	Config          config.Config
	DB              sqlx.SqlConn
	AdminUserModel  model.AdminUserModel
	ShortLinkModel  model.ShortLinkModel
	VisitEventModel model.VisitEventModel
	DailyStatModel  model.LinkDailyStatModel
}

// NewServiceContext creates a link-rpc service context.
func NewServiceContext(c config.Config) *ServiceContext {
	db := sqlx.NewMysql(c.Dependencies.MySQL.DataSource)

	return &ServiceContext{
		Config:          c,
		DB:              db,
		AdminUserModel:  model.NewAdminUserModel(db, c.CacheRedis),
		ShortLinkModel:  model.NewShortLinkModel(db, c.CacheRedis),
		VisitEventModel: model.NewVisitEventModel(db),
		DailyStatModel:  model.NewLinkDailyStatModel(db),
	}
}
