package cleanup

import (
	"context"
	"fmt"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/metrics"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"

	"github.com/zeromicro/go-zero/core/logx"
)

const shortLinkColumns = "`id`,`code`,`origin_url`,`title`,`description`," +
	"`status`,`expire_at`,`created_by`,`created_at`,`updated_at`,`deleted_at`"

func (r *Runner) cleanArchivedLinks(ctx context.Context) {
	cutoff := time.Now().UTC().AddDate(0, 0, -r.cfg.ShortLinkRetentionDays)
	total := int64(0)
	for {
		var rows []*model.ShortLink
		fetchQuery := fmt.Sprintf(
			"select "+shortLinkColumns+
				" from `short_link` where `deleted_at` is not null and `deleted_at` < ? limit %d",
			r.cfg.CleanupBatchSize,
		)
		if err := r.db.QueryRowsCtx(ctx, &rows, fetchQuery, cutoff); err != nil {
			logx.Errorw("cleanup: short_link fetch failed", logx.Field("error", err.Error()))
			return
		}
		if len(rows) == 0 {
			break
		}
		for _, link := range rows {
			if err := r.archiveLink(ctx, link); err != nil {
				logx.Errorw("cleanup: archive link failed",
					logx.Field("id", link.Id),
					logx.Field("code", link.Code),
					logx.Field("error", err.Error()),
				)
				continue
			}
			metrics.CleanupDeletedRowsTotal.Inc("short_link")
			total++
		}
		if len(rows) < r.cfg.CleanupBatchSize {
			break
		}
	}
	logx.Infow("cleanup: short_link archive done", logx.Field("archived", total))
}

func (r *Runner) archiveLink(ctx context.Context, link *model.ShortLink) error {
	// 1. Insert into archive — idempotent on re-run.
	archiveQuery := "insert ignore into `short_link_archive` (" + shortLinkColumns + ")" +
		" values (?,?,?,?,?,?,?,?,?,?,?)"
	if _, err := r.db.ExecCtx(ctx, archiveQuery,
		link.Id, link.Code, link.OriginUrl, link.Title, link.Description,
		link.Status, link.ExpireAt, link.CreatedBy, link.CreatedAt, link.UpdatedAt, link.DeletedAt,
	); err != nil {
		return fmt.Errorf("insert archive: %w", err)
	}
	// 2. Reserve the code permanently so it can never be recreated.
	if _, err := r.db.ExecCtx(ctx,
		"insert ignore into `reserved_code` (`code`, `reserved_at`) values (?, now())", link.Code,
	); err != nil {
		return fmt.Errorf("insert reserved_code: %w", err)
	}
	// 3. Hard-delete from short_link.
	if _, err := r.db.ExecCtx(ctx,
		"delete from `short_link` where `id` = ?", link.Id,
	); err != nil {
		return fmt.Errorf("delete short_link: %w", err)
	}
	return nil
}
