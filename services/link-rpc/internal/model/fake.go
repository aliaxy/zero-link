//go:build testing

package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// FakeShortLinkModel is a test double for ShortLinkModel.
type FakeShortLinkModel struct {
	Link *ShortLink
	Err  error
}

func (f FakeShortLinkModel) Insert(_ context.Context, _ *ShortLink) (sql.Result, error) {
	return nil, nil
}

func (f FakeShortLinkModel) FindOne(_ context.Context, _ int64) (*ShortLink, error) {
	return f.Link, f.Err
}

func (f FakeShortLinkModel) FindOneByCode(_ context.Context, _ string) (*ShortLink, error) {
	return f.Link, f.Err
}

func (f FakeShortLinkModel) FindOneNotDeleted(_ context.Context, _ int64) (*ShortLink, error) {
	return f.Link, f.Err
}
func (f FakeShortLinkModel) Update(_ context.Context, _ *ShortLink) error { return nil }
func (f FakeShortLinkModel) Delete(_ context.Context, _ int64) error      { return nil }
func (f FakeShortLinkModel) List(_ context.Context, _ ShortLinkListFilter) ([]*ShortLink, int64, error) {
	return nil, 0, nil
}
func (f FakeShortLinkModel) SoftDelete(_ context.Context, _ int64, _ time.Time) error { return nil }
func (f FakeShortLinkModel) withSession(_ sqlx.Session) ShortLinkModel                { return f }
