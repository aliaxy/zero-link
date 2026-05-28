// Package model contains generated database models for link-rpc.
package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ShortLinkModel = (*customShortLinkModel)(nil)

type (
	// ShortLinkModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShortLinkModel.
	ShortLinkModel interface {
		shortLinkModel
		FindOneNotDeleted(ctx context.Context, id int64) (*ShortLink, error)
		List(ctx context.Context, filter ShortLinkListFilter) ([]*ShortLink, int64, error)
		SoftDelete(ctx context.Context, id int64, deletedAt time.Time) error
		withSession(session sqlx.Session) ShortLinkModel
	}

	// ShortLinkListFilter contains management short-link list filters.
	ShortLinkListFilter struct {
		Page     int64
		PageSize int64
		Status   int64
		Keyword  string
	}

	customShortLinkModel struct {
		*defaultShortLinkModel
	}
)

// NewShortLinkModel returns a model for the database table.
func NewShortLinkModel(conn sqlx.SqlConn) ShortLinkModel {
	return &customShortLinkModel{
		defaultShortLinkModel: newShortLinkModel(conn),
	}
}

func (m *customShortLinkModel) withSession(session sqlx.Session) ShortLinkModel {
	return NewShortLinkModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customShortLinkModel) FindOneNotDeleted(ctx context.Context, id int64) (*ShortLink, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `deleted_at` is null limit 1", shortLinkRows, m.table)
	var resp ShortLink
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customShortLinkModel) List(ctx context.Context, filter ShortLinkListFilter) ([]*ShortLink, int64, error) {
	where, args := buildShortLinkListWhere(filter)

	countQuery := fmt.Sprintf("select count(*) from %s %s", m.table, where)
	var total int64
	if err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)
	query := fmt.Sprintf(
		"select %s from %s %s order by `created_at` desc, `id` desc limit ? offset ?",
		shortLinkRows,
		m.table,
		where,
	)

	var links []*ShortLink
	if err := m.conn.QueryRowsCtx(ctx, &links, query, args...); err != nil {
		return nil, 0, err
	}

	return links, total, nil
}

func (m *customShortLinkModel) SoftDelete(ctx context.Context, id int64, deletedAt time.Time) error {
	query := fmt.Sprintf("update %s set `deleted_at` = ? where `id` = ? and `deleted_at` is null", m.table)
	result, err := m.conn.ExecCtx(ctx, query, sql.NullTime{Time: deletedAt, Valid: true}, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func buildShortLinkListWhere(filter ShortLinkListFilter) (string, []any) {
	conditions := []string{"`deleted_at` is null"}
	args := make([]any, 0, 4)

	if filter.Status > 0 {
		conditions = append(conditions, "`status` = ?")
		args = append(args, filter.Status)
	}
	if filter.Keyword != "" {
		conditions = append(conditions, "(`code` like ? or `title` like ? or `origin_url` like ?)")
		keyword := "%" + filter.Keyword + "%"
		args = append(args, keyword, keyword, keyword)
	}

	return "where " + strings.Join(conditions, " and "), args
}
