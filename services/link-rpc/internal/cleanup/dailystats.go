package cleanup

import (
	"context"
	"fmt"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/metrics"

	"github.com/zeromicro/go-zero/core/logx"
)

func (r *Runner) cleanDailyStats(ctx context.Context) {
	cutoff := time.Now().UTC().AddDate(0, 0, -r.cfg.DailyStatRetentionDays).Format("2006-01-02")
	total := int64(0)
	for {
		result, err := r.db.ExecCtx(ctx,
			fmt.Sprintf("delete from `link_daily_stat` where `stat_date` < ? limit %d", r.cfg.CleanupBatchSize),
			cutoff,
		)
		if err != nil {
			logx.Errorw("cleanup: link_daily_stat delete failed", logx.Field("error", err.Error()))
			return
		}
		affected, err := result.RowsAffected()
		if err != nil {
			logx.Errorw("cleanup: link_daily_stat rows affected failed", logx.Field("error", err.Error()))
			return
		}
		if affected == 0 {
			break
		}
		metrics.CleanupDeletedRowsTotal.Add(float64(affected), "link_daily_stat")
		total += affected
	}
	logx.Infow("cleanup: link_daily_stat done", logx.Field("deleted", total))
}
