package logic

import (
	"context"
	"testing"
	"time"

	"github.com/aliaxy/zero-link/services/link-rpc/internal/model"
	"github.com/aliaxy/zero-link/services/link-rpc/internal/svc"
	linkv1 "github.com/aliaxy/zero-link/services/link-rpc/pb/link/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeDailyStatModelWithRange struct {
	model.LinkDailyStatModel
	rows []*model.LinkDailyStat
	err  error
}

func (f fakeDailyStatModelWithRange) FindByLinkIDAndDateRange(
	_ context.Context, _ int64, _, _ time.Time,
) ([]*model.LinkDailyStat, error) {
	return f.rows, f.err
}

func newGetLinkStatsLogic(m model.LinkDailyStatModel) *GetLinkStatsLogic {
	return NewGetLinkStatsLogic(context.Background(), &svc.ServiceContext{
		DailyStatModel: m,
	})
}

func TestGetLinkStats_NoData(t *testing.T) {
	logic := newGetLinkStatsLogic(fakeDailyStatModelWithRange{rows: nil})
	resp, err := logic.GetLinkStats(&linkv1.GetLinkStatsRequest{
		LinkId: 1,
		From:   "2026-05-01",
		To:     "2026-05-29",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 0 {
		t.Fatalf("want 0 items, got %d", len(resp.Items))
	}
}

func TestGetLinkStats_WithData(t *testing.T) {
	statDate, _ := time.Parse("2006-01-02", "2026-05-15")
	logic := newGetLinkStatsLogic(fakeDailyStatModelWithRange{
		rows: []*model.LinkDailyStat{
			{LinkId: 1, StatDate: statDate, Pv: 42, Uv: 7},
		},
	})
	resp, err := logic.GetLinkStats(&linkv1.GetLinkStatsRequest{
		LinkId: 1,
		From:   "2026-05-01",
		To:     "2026-05-29",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("want 1 item, got %d", len(resp.Items))
	}
	if resp.Items[0].StatDate != "2026-05-15" {
		t.Fatalf("stat_date = %q, want 2026-05-15", resp.Items[0].StatDate)
	}
	if resp.Items[0].Pv != 42 {
		t.Fatalf("pv = %d, want 42", resp.Items[0].Pv)
	}
}

func TestGetLinkStats_InvalidDateRange(t *testing.T) {
	logic := newGetLinkStatsLogic(fakeDailyStatModelWithRange{})
	_, err := logic.GetLinkStats(&linkv1.GetLinkStatsRequest{
		LinkId: 1,
		From:   "2026-05-29",
		To:     "2026-05-01",
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("want InvalidArgument, got %v", err)
	}
}
