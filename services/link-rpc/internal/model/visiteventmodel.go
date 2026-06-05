package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ VisitEventModel = (*customVisitEventModel)(nil)

type (
	// VisitEventModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVisitEventModel.
	VisitEventModel interface {
		visitEventModel
		withSession(session sqlx.Session) VisitEventModel
		HasVisitedToday(ctx context.Context, linkID int64, ipHash, statDate string) (bool, error)
	}

	customVisitEventModel struct {
		*defaultVisitEventModel
	}
)

// NewVisitEventModel returns a model for the database table.
func NewVisitEventModel(conn sqlx.SqlConn) VisitEventModel {
	return &customVisitEventModel{
		defaultVisitEventModel: newVisitEventModel(conn),
	}
}

func (m *customVisitEventModel) withSession(session sqlx.Session) VisitEventModel {
	return NewVisitEventModel(sqlx.NewSqlConnFromSession(session))
}

// HasVisitedToday reports whether an ip_hash has already visited the given link on statDate.
func (m *customVisitEventModel) HasVisitedToday(ctx context.Context, linkID int64, ipHash, statDate string) (bool, error) {
	query := fmt.Sprintf(
		"select count(1) from %s where `link_id` = ? and `ip_hash` = ? and date(`visited_at`) = ?",
		m.table,
	)
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, linkID, ipHash, statDate)
	return count > 0, err
}
