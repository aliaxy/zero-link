package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LinkDailyStatModel = (*customLinkDailyStatModel)(nil)

type (
	// LinkDailyStatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLinkDailyStatModel.
	LinkDailyStatModel interface {
		linkDailyStatModel
		withSession(session sqlx.Session) LinkDailyStatModel
		UpsertStats(ctx context.Context, linkID int64, statDate string, isNewUV bool) error
		FindByLinkIDAndDateRange(ctx context.Context, linkID int64, from, to time.Time) ([]*LinkDailyStat, error)
	}

	customLinkDailyStatModel struct {
		*defaultLinkDailyStatModel
	}
)

// NewLinkDailyStatModel returns a model for the database table.
func NewLinkDailyStatModel(conn sqlx.SqlConn) LinkDailyStatModel {
	return &customLinkDailyStatModel{
		defaultLinkDailyStatModel: newLinkDailyStatModel(conn),
	}
}

func (m *customLinkDailyStatModel) withSession(session sqlx.Session) LinkDailyStatModel {
	return NewLinkDailyStatModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customLinkDailyStatModel) UpsertStats(ctx context.Context, linkID int64, statDate string, isNewUV bool) error {
	var query string
	if isNewUV {
		query = fmt.Sprintf(
			"insert into %s (`link_id`, `stat_date`, `pv`, `uv`) values (?, ?, 1, 1)"+
				" on duplicate key update `pv` = `pv` + 1, `uv` = `uv` + 1",
			m.table,
		)
	} else {
		query = fmt.Sprintf(
			"insert into %s (`link_id`, `stat_date`, `pv`, `uv`) values (?, ?, 1, 0) on duplicate key update `pv` = `pv` + 1",
			m.table,
		)
	}
	_, err := m.conn.ExecCtx(ctx, query, linkID, statDate)
	return err
}

func (m *customLinkDailyStatModel) FindByLinkIDAndDateRange(
	ctx context.Context, linkID int64, from, to time.Time,
) ([]*LinkDailyStat, error) {
	query := fmt.Sprintf(
		"select %s from %s where `link_id` = ? and `stat_date` between ? and ? order by `stat_date` asc",
		linkDailyStatRows, m.table,
	)
	var rows []*LinkDailyStat
	if err := m.conn.QueryRowsCtx(ctx, &rows, query, linkID, from, to); err != nil {
		return nil, err
	}
	return rows, nil
}
