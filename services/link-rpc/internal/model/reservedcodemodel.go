package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ReservedCodeModel = (*customReservedCodeModel)(nil)

type (
	// ReservedCodeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customReservedCodeModel.
	ReservedCodeModel interface {
		reservedCodeModel
		withSession(session sqlx.Session) ReservedCodeModel
		Exists(ctx context.Context, code string) (bool, error)
	}

	customReservedCodeModel struct {
		*defaultReservedCodeModel
	}
)

// NewReservedCodeModel returns a model for the database table.
func NewReservedCodeModel(conn sqlx.SqlConn) ReservedCodeModel {
	return &customReservedCodeModel{
		defaultReservedCodeModel: newReservedCodeModel(conn),
	}
}

func (m *customReservedCodeModel) withSession(session sqlx.Session) ReservedCodeModel {
	return NewReservedCodeModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customReservedCodeModel) Exists(ctx context.Context, code string) (bool, error) {
	var count int64
	query := "select count(*) from `reserved_code` where `code` = ? limit 1"
	err := m.conn.QueryRowCtx(ctx, &count, query, code)
	return count > 0, err
}
