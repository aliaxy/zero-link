// Package model contains generated database models for link-rpc.
package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
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
		// HardDelete permanently removes a short link row and invalidates its Redis cache entries.
		// Used by the cleanup runner after the archive transaction commits.
		HardDelete(ctx context.Context, id int64, code string) error
		// ArchiveInsertIgnore copies a short link row into short_link_archive.
		// Idempotent: duplicate ids are silently skipped (INSERT IGNORE).
		ArchiveInsertIgnore(ctx context.Context, link *ShortLink) error
		// ListSoftDeletedBefore returns at most limit soft-deleted rows whose
		// deleted_at is earlier than cutoff. Used by the cleanup runner.
		ListSoftDeletedBefore(ctx context.Context, cutoff time.Time, limit int) ([]*ShortLink, error)
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
		cacheConf cache.CacheConf
	}
)

// NewShortLinkModel returns a model for the database table.
func NewShortLinkModel(conn sqlx.SqlConn, c cache.CacheConf) ShortLinkModel {
	return &customShortLinkModel{
		defaultShortLinkModel: newShortLinkModel(conn, c),
		cacheConf:             c,
	}
}

func (m *customShortLinkModel) withSession(session sqlx.Session) ShortLinkModel {
	return NewShortLinkModel(sqlx.NewSqlConnFromSession(session), m.cacheConf)
}

func (m *customShortLinkModel) FindOneNotDeleted(ctx context.Context, id int64) (*ShortLink, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? and `deleted_at` is null limit 1", shortLinkRows, m.table)
	var resp ShortLink
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, id)
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
	if err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...); err != nil {
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
	if err := m.QueryRowsNoCacheCtx(ctx, &links, query, args...); err != nil {
		return nil, 0, err
	}

	return links, total, nil
}

func (m *customShortLinkModel) SoftDelete(ctx context.Context, id int64, deletedAt time.Time) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	idKey := fmt.Sprintf("%s%v", cacheShortLinkIdPrefix, id)
	codeKey := fmt.Sprintf("%s%v", cacheShortLinkCodePrefix, data.Code)

	result, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("update %s set `deleted_at` = ? where `id` = ? and `deleted_at` is null", m.table)
		return conn.ExecCtx(ctx, query, sql.NullTime{Time: deletedAt, Valid: true}, id)
	}, idKey, codeKey)
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

func (m *customShortLinkModel) HardDelete(ctx context.Context, id int64, code string) error {
	idKey := fmt.Sprintf("%s%v", cacheShortLinkIdPrefix, id)
	codeKey := fmt.Sprintf("%s%v", cacheShortLinkCodePrefix, code)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, idKey, codeKey)
	return err
}

func (m *customShortLinkModel) ArchiveInsertIgnore(ctx context.Context, link *ShortLink) error {
	const archiveTable = "`short_link_archive`"
	const columns = "`id`,`code`,`origin_url`,`title`,`description`," +
		"`status`,`expire_at`,`created_by`,`created_at`,`updated_at`,`deleted_at`"
	query := fmt.Sprintf("insert ignore into %s (%s) values (?,?,?,?,?,?,?,?,?,?,?)", archiveTable, columns)
	// ExecCtx with no cache keys: uses the underlying connection (session-aware)
	// without touching Redis, since the archive table has no cache layer.
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, query,
			link.Id, link.Code, link.OriginUrl, link.Title, link.Description,
			link.Status, link.ExpireAt, link.CreatedBy, link.CreatedAt, link.UpdatedAt, link.DeletedAt,
		)
	})
	return err
}

func (m *customShortLinkModel) ListSoftDeletedBefore(
	ctx context.Context, cutoff time.Time, limit int,
) ([]*ShortLink, error) {
	query := fmt.Sprintf(
		"select %s from %s where `deleted_at` is not null and `deleted_at` < ? limit ?",
		shortLinkRows, m.table,
	)
	var rows []*ShortLink
	if err := m.QueryRowsNoCacheCtx(ctx, &rows, query, cutoff, limit); err != nil {
		return nil, err
	}
	return rows, nil
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
