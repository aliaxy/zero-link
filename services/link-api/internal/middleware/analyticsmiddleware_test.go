package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aliaxy/zero-link/services/link-api/pkg/httputil"
	"github.com/aliaxy/zero-link/services/link-rpc/client/analyticsservice"
	"go.uber.org/goleak"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	// IgnoreCurrent skips goroutines started by go-zero's init() functions
	// (proc signal handler, stat usage collector) that are not under our control.
	goleak.VerifyTestMain(m, goleak.IgnoreCurrent())
}

// recordingLinkService captures RecordVisit calls via a channel for deterministic tests.
type recordingLinkService struct {
	analyticsservice.AnalyticsService
	visitCh chan *analyticsservice.RecordVisitRequest
}

func (r *recordingLinkService) RecordVisit(
	_ context.Context,
	req *analyticsservice.RecordVisitRequest,
	_ ...grpc.CallOption,
) (*analyticsservice.RecordVisitResponse, error) {
	r.visitCh <- req
	return &analyticsservice.RecordVisitResponse{}, nil
}

func newRecordingMiddleware() (*AnalyticsMiddleware, chan *analyticsservice.RecordVisitRequest) {
	ch := make(chan *analyticsservice.RecordVisitRequest, 1)
	return NewAnalyticsMiddleware(&recordingLinkService{visitCh: ch}), ch
}

func TestAnalyticsMiddleware_302TriggersRecordVisit(t *testing.T) {
	mw, ch := newRecordingMiddleware()

	handler := mw.Handle(func(w http.ResponseWriter, _ *http.Request) {
		http.Redirect(w, httptest.NewRequest(http.MethodGet, "/abc123", nil), "https://example.com", http.StatusFound)
	})

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	handler(httptest.NewRecorder(), req)

	select {
	case got := <-ch:
		if got.Code != "abc123" {
			t.Fatalf("code = %q, want abc123", got.Code)
		}
		if got.Ip != "1.2.3.4" {
			t.Fatalf("ip = %q, want 1.2.3.4", got.Ip)
		}
		if got.UserAgent != "Mozilla/5.0" {
			t.Fatalf("user_agent = %q, want Mozilla/5.0", got.UserAgent)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for RecordVisit call")
	}
}

func TestAnalyticsMiddleware_404NoRecordVisit(t *testing.T) {
	mw, ch := newRecordingMiddleware()

	handler := mw.Handle(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	handler(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/missing", nil))

	select {
	case <-ch:
		t.Fatal("RecordVisit should not be called on 404")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestAnalyticsMiddleware_403NoRecordVisit(t *testing.T) {
	mw, ch := newRecordingMiddleware()

	handler := mw.Handle(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	handler(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/gone", nil))

	select {
	case <-ch:
		t.Fatal("RecordVisit should not be called on 403")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestExtractIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "1.2.3.4, 10.0.0.1")
	if got := httputil.ExtractClientIP(req); got != "1.2.3.4" {
		t.Fatalf("extractIP = %q, want 1.2.3.4", got)
	}
}

func TestExtractIP_FallbackRemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "5.6.7.8:9999"
	if got := httputil.ExtractClientIP(req); got != "5.6.7.8" {
		t.Fatalf("extractIP = %q, want 5.6.7.8", got)
	}
}
