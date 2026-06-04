package cleanup

import (
	"context"
	"fmt"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/metrics"

	"github.com/zeromicro/go-zero/core/logx"
)

func (r *Runner) cleanVisitEvents(ctx context.Context) {
	cutoff := time.Now().UTC().AddDate(0, 0, -r.cfg.VisitEventRetentionDays)
	total := int64(0)
	for {
		result, err := r.db.ExecCtx(ctx,
			fmt.Sprintf("delete from `visit_event` where `visited_at` < ? limit %d", r.cfg.CleanupBatchSize),
			cutoff,
		)
		if err != nil {
			logx.Errorw("cleanup: visit_event delete failed", logx.Field("error", err.Error()))
			return
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			break
		}
		for range affected {
			metrics.CleanupDeletedRowsTotal.Inc("visit_event")
		}
		total += affected
	}
	logx.Infow("cleanup: visit_event done", logx.Field("deleted", total))
}
