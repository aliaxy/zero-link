// Package cleanup runs periodic data-retention jobs for link-rpc.
package cleanup

import (
	"context"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// Runner executes periodic data retention cleanup jobs.
type Runner struct {
	db  sqlx.SqlConn
	cfg config.RetentionConfig
}

// NewRunner creates a Runner.
func NewRunner(db sqlx.SqlConn, cfg config.RetentionConfig) *Runner {
	return &Runner{db: db, cfg: cfg}
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
}
