package cleanup_test

import (
	"context"
	"testing"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/cleanup"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/config"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m, goleak.IgnoreCurrent())
}

func TestRunner_StartCancellation(_ *testing.T) {
	r := cleanup.NewRunner(nil, nil, nil, config.RetentionConfig{
		VisitEventRetentionDays: 90,
		ShortLinkRetentionDays:  365,
		DailyStatRetentionDays:  730,
		CleanupBatchSize:        1000,
	})
	ctx, cancel := context.WithCancel(context.Background())
	r.Start(ctx)
	time.Sleep(10 * time.Millisecond)
	cancel()
	time.Sleep(20 * time.Millisecond)
}
