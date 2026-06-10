package cleanup

import (
	"context"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/metrics"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"

	"github.com/zeromicro/go-zero/core/logx"
)

func (r *Runner) cleanArchivedLinks(ctx context.Context) {
	cutoff := time.Now().UTC().AddDate(0, 0, -r.cfg.ShortLinkRetentionDays)
	total := int64(0)
	for {
		rows, err := r.shortLinkModel.ListSoftDeletedBefore(ctx, cutoff, r.cfg.CleanupBatchSize)
		if err != nil {
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
	// Phase 1 (atomic): copy row to short_link_archive + insert into reserved_code.
	// If the process crashes after this but before phase 2, the next cleanup run
	// retries: both inserts are INSERT IGNORE (idempotent), then phase 2 succeeds.
	if err := r.archiver.ArchiveAndReserveCode(ctx, link); err != nil {
		return err
	}
	// Phase 2: hard-delete from short_link via the model layer so the Redis cache
	// entries for this id/code are evicted immediately.
	return r.shortLinkModel.HardDelete(ctx, link.Id, link.Code)
}
