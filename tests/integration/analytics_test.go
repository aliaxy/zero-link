//go:build integration

package integration_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestAnalytics_RecordVisitAndStats(t *testing.T) {
	token := doLogin(t, testAdminUsername, testAdminPassword)
	originURL := "https://example.com/analytics-" + uniqueName(t)
	id, code := createLink(t, token, originURL)

	// Register cleanups in LIFO order: analytics rows first, then the link.
	t.Cleanup(func() {
		dbExec(t, "DELETE FROM visit_event WHERE link_id = ?", id)
		dbExec(t, "DELETE FROM link_daily_stat WHERE link_id = ?", id)
	})
	t.Cleanup(func() { deleteLink(t, token, id) })

	// Trigger redirect → AnalyticsMiddleware fires an async goroutine → RecordVisit gRPC.
	redirectResp := doRequestNoRedirect(t, http.MethodGet, "/"+code, nil, "")
	redirectResp.Body.Close()
	assertStatus(t, redirectResp, http.StatusFound)

	// Poll stats until PV ≥ 1 or 8s deadline.
	// The AnalyticsMiddleware goroutine has a 3s context timeout; 8s gives
	// room for go compilation overhead and slow CI machines.
	today := time.Now().UTC().Format("2006-01-02")
	statsURL := fmt.Sprintf("/admin/links/%d/stats?from=%s&to=%s", id, today, today)

	var lastPV int64
	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		resp := doRequest(t, http.MethodGet, statsURL, nil, token)
		assertStatus(t, resp, http.StatusOK)
		var envelope struct {
			Data struct {
				Items []struct {
					StatDate string `json:"stat_date"`
					Pv       int64  `json:"pv"`
				} `json:"items"`
			} `json:"data"`
		}
		decodeJSON(t, resp, &envelope)
		if len(envelope.Data.Items) > 0 && envelope.Data.Items[0].Pv >= 1 {
			lastPV = envelope.Data.Items[0].Pv
			if envelope.Data.Items[0].StatDate != today {
				t.Fatalf("stat_date = %q, want %q", envelope.Data.Items[0].StatDate, today)
			}
			break
		}
		time.Sleep(300 * time.Millisecond)
	}

	if lastPV < 1 {
		t.Fatalf("expected pv >= 1 after redirect, got %d (analytics may not have flushed)", lastPV)
	}
}

func TestAnalytics_MultipleRedirects(t *testing.T) {
	token := doLogin(t, testAdminUsername, testAdminPassword)
	originURL := "https://example.com/analytics-multi-" + uniqueName(t)
	id, code := createLink(t, token, originURL)

	t.Cleanup(func() {
		dbExec(t, "DELETE FROM visit_event WHERE link_id = ?", id)
		dbExec(t, "DELETE FROM link_daily_stat WHERE link_id = ?", id)
	})
	t.Cleanup(func() { deleteLink(t, token, id) })

	// Fire 3 redirects.
	for range 3 {
		resp := doRequestNoRedirect(t, http.MethodGet, "/"+code, nil, "")
		resp.Body.Close()
	}

	today := time.Now().UTC().Format("2006-01-02")
	statsURL := fmt.Sprintf("/admin/links/%d/stats?from=%s&to=%s", id, today, today)

	var lastPV int64
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		resp := doRequest(t, http.MethodGet, statsURL, nil, token)
		assertStatus(t, resp, http.StatusOK)
		var envelope struct {
			Data struct {
				Items []struct {
					Pv int64 `json:"pv"`
				} `json:"items"`
			} `json:"data"`
		}
		decodeJSON(t, resp, &envelope)
		if len(envelope.Data.Items) > 0 {
			lastPV = envelope.Data.Items[0].Pv
			if lastPV >= 3 {
				break
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	if lastPV < 3 {
		t.Fatalf("expected pv >= 3 after 3 redirects, got %d", lastPV)
	}
}

func TestAnalytics_NoStatsForNewLink(t *testing.T) {
	token := doLogin(t, testAdminUsername, testAdminPassword)
	id, _ := createLink(t, token, "https://example.com/no-stats-"+uniqueName(t))
	t.Cleanup(func() { deleteLink(t, token, id) })

	today := time.Now().UTC().Format("2006-01-02")
	statsURL := fmt.Sprintf("/admin/links/%d/stats?from=%s&to=%s", id, today, today)

	resp := doRequest(t, http.MethodGet, statsURL, nil, token)
	assertStatus(t, resp, http.StatusOK)

	var envelope struct {
		Data struct {
			Items []any `json:"items"`
		} `json:"data"`
	}
	decodeJSON(t, resp, &envelope)
	if len(envelope.Data.Items) != 0 {
		t.Fatalf("expected 0 stat items for new link, got %d", len(envelope.Data.Items))
	}
}
