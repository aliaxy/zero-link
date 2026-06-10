package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// LinkArchiver coordinates the atomic archival of a soft-deleted short link.
// It owns the transaction that spans ShortLinkModel and ReservedCodeModel,
// keeping withSession unexported and transaction logic inside the model layer.
type LinkArchiver struct {
	shortLinkModel    ShortLinkModel
	reservedCodeModel ReservedCodeModel
	db                sqlx.SqlConn
}

// NewLinkArchiver creates a LinkArchiver.
func NewLinkArchiver(db sqlx.SqlConn, slm ShortLinkModel, rcm ReservedCodeModel) *LinkArchiver {
	return &LinkArchiver{shortLinkModel: slm, reservedCodeModel: rcm, db: db}
}

// ArchiveAndReserveCode atomically copies link into short_link_archive and
// inserts its code into reserved_code, preventing the code from ever being
// reused after the link is permanently removed from short_link.
func (a *LinkArchiver) ArchiveAndReserveCode(ctx context.Context, link *ShortLink) error {
	return a.db.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := a.shortLinkModel.withSession(session).ArchiveInsertIgnore(ctx, link); err != nil {
			return err
		}
		return a.reservedCodeModel.withSession(session).Reserve(ctx, link.Code)
	})
}
