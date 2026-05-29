package logic

import (
	"context"
	"errors"
	"testing"

	"github.com/aliaxy/zero-link/services/link-api/internal/svc"
	"github.com/aliaxy/zero-link/services/link-api/internal/types"
	"github.com/aliaxy/zero-link/services/link-rpc/linkservice"
)

func newGetLinkStatsAPILogic(svc *svc.ServiceContext) *GetLinkStatsLogic {
	return NewGetLinkStatsLogic(context.Background(), svc)
}

func TestGetLinkStatsLogic_Success(t *testing.T) {
	logic := newGetLinkStatsAPILogic(&svc.ServiceContext{
		LinkRPC: fakeLinkService{
			statsResp: &linkservice.GetLinkStatsResponse{
				LinkId: 1,
				Items: []*linkservice.DailyStat{
					{StatDate: "2026-05-15", Pv: 42, Uv: 7},
				},
			},
		},
	})

	resp, err := logic.GetLinkStats(&types.GetLinkStatsRequest{
		Id:   1,
		From: "2026-05-01",
		To:   "2026-05-29",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != "SUCCESS" {
		t.Fatalf("code = %q, want SUCCESS", resp.Code)
	}
	if resp.Data.LinkId != 1 {
		t.Fatalf("link_id = %d, want 1", resp.Data.LinkId)
	}
	if len(resp.Data.Items) != 1 {
		t.Fatalf("items length = %d, want 1", len(resp.Data.Items))
	}
	if resp.Data.Items[0].StatDate != "2026-05-15" {
		t.Fatalf("stat_date = %q, want 2026-05-15", resp.Data.Items[0].StatDate)
	}
	if resp.Data.Items[0].Pv != 42 {
		t.Fatalf("pv = %d, want 42", resp.Data.Items[0].Pv)
	}
}

func TestGetLinkStatsLogic_RPCError(t *testing.T) {
	logic := newGetLinkStatsAPILogic(&svc.ServiceContext{
		LinkRPC: fakeLinkService{
			err: errors.New("rpc failure"),
		},
	})

	_, err := logic.GetLinkStats(&types.GetLinkStatsRequest{Id: 1})
	if err == nil {
		t.Fatal("expected error from RPC failure, got nil")
	}
}
