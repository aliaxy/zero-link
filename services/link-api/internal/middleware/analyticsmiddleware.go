// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/aliaxy/zero-link/services/link-api/pkg/httputil"
	"github.com/aliaxy/zero-link/services/link-rpc/client/analyticsservice"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	analyticsWorkers   = 8
	analyticsQueueSize = 2000
)

type analyticsJob struct {
	req *analyticsservice.RecordVisitRequest
	log logx.Logger
}

// AnalyticsMiddleware fires a non-blocking RecordVisit RPC after every 302 redirect.
// A bounded channel-backed worker pool absorbs burst traffic and drops events when full
// rather than spawning unbounded goroutines.
type AnalyticsMiddleware struct {
	linkRPC        analyticsservice.AnalyticsService
	trustedProxies []string
	jobs           chan *analyticsJob
}

// NewAnalyticsMiddleware creates an AnalyticsMiddleware and starts its worker pool.
func NewAnalyticsMiddleware(linkRPC analyticsservice.AnalyticsService, trustedProxies []string) *AnalyticsMiddleware {
	m := &AnalyticsMiddleware{
		linkRPC:        linkRPC,
		trustedProxies: trustedProxies,
		jobs:           make(chan *analyticsJob, analyticsQueueSize),
	}
	for range analyticsWorkers {
		go m.runWorker()
	}
	return m
}

func (m *AnalyticsMiddleware) runWorker() {
	for job := range m.jobs {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		if _, err := m.linkRPC.RecordVisit(ctx, job.req); err != nil {
			job.log.Errorf("record visit async: %v", err)
		}
		cancel()
	}
}

// Stop drains the worker pool. Call this once when the service shuts down.
func (m *AnalyticsMiddleware) Stop() {
	close(m.jobs)
}

// Handle wraps the redirect handler and records visit events asynchronously.
func (m *AnalyticsMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(rw, r)

		if rw.status != http.StatusFound {
			return
		}

		job := &analyticsJob{
			req: &analyticsservice.RecordVisitRequest{
				Code:      strings.TrimPrefix(r.URL.Path, "/"),
				VisitedAt: time.Now().UTC().Format(time.RFC3339Nano),
				Ip:        httputil.ExtractClientIP(r, m.trustedProxies),
				UserAgent: r.Header.Get("User-Agent"),
				Referer:   r.Header.Get("Referer"),
			},
			log: logx.WithContext(r.Context()),
		}

		select {
		case m.jobs <- job:
		default:
			job.log.Errorf("analytics queue full, dropping visit for code %s", job.req.Code)
		}
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
