// Package model contains generated database models for link-rpc.
package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminUserModel = (*customAdminUserModel)(nil)

type (
	// AdminUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminUserModel.
	AdminUserModel interface {
		adminUserModel
		withSession(session sqlx.Session) AdminUserModel
	}

	customAdminUserModel struct {
		*defaultAdminUserModel
		cacheConf cache.CacheConf
	}
)

// NewAdminUserModel returns a model for the database table.
func NewAdminUserModel(conn sqlx.SqlConn, c cache.CacheConf) AdminUserModel {
	return &customAdminUserModel{
		defaultAdminUserModel: newAdminUserModel(conn, c),
		cacheConf:             c,
	}
}

func (m *customAdminUserModel) withSession(session sqlx.Session) AdminUserModel {
	return NewAdminUserModel(sqlx.NewSqlConnFromSession(session), m.cacheConf)
}
