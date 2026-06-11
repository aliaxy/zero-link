// Package cleanup runs periodic data-retention jobs for link-rpc.
package cleanup

import (
	"context"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/config"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/metrics"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/pkg/filter"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// Runner executes periodic data retention cleanup jobs.
type Runner struct {
	db             sqlx.SqlConn
	shortLinkModel model.ShortLinkModel
	archiver       *model.LinkArchiver
	cfg            config.RetentionConfig
	filter         *filter.CodeFilter
	filterCap      uint
}

// NewRunner creates a Runner.
// shortLinkModel is used for listing soft-deleted links and for the hard-delete
// step (which also invalidates the Redis cache).
// archiver coordinates the atomic archive+reserve transaction.
// Pass nil for either to skip that step (useful in tests).
// cf and filterCap are optional (pass nil/0 to skip fill-ratio reporting).
func NewRunner(
	db sqlx.SqlConn,
	shortLinkModel model.ShortLinkModel,
	archiver *model.LinkArchiver,
	cfg config.RetentionConfig,
	cf *filter.CodeFilter,
	filterCap uint,
) *Runner {
	return &Runner{
		db:             db,
		shortLinkModel: shortLinkModel,
		archiver:       archiver,
		cfg:            cfg,
		filter:         cf,
		filterCap:      filterCap,
	}
}

// Start launches a background goroutine that runs cleanup every 24 hours.
// The goroutine stops when ctx is cancelled.
func (r *Runner) Start(ctx context.Context) {
	go func() {
		r.runOnce(ctx)
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				r.runOnce(ctx)
			}
		}
	}()
}

func (r *Runner) runOnce(ctx context.Context) {
	if r.db == nil {
		return
	}
	logx.Info("cleanup: starting retention run")
	r.cleanVisitEvents(ctx)
	r.cleanArchivedLinks(ctx)
	r.cleanDailyStats(ctx)
	logx.Info("cleanup: retention run complete")
	r.reportFilterFillRatio()
}

func (r *Runner) reportFilterFillRatio() {
	if r.filter == nil || r.filterCap == 0 {
		return
	}
	metrics.FilterFillRatio.Set(float64(r.filter.Count()) / float64(r.filterCap))
}
