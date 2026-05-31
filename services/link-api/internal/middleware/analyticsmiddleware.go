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

// AnalyticsMiddleware fires a non-blocking RecordVisit RPC after every 302 redirect.
// Known limitation: goroutines are unbounded under high traffic; Stage 6 will replace
// with a channel-based worker pool.
type AnalyticsMiddleware struct {
	linkRPC analyticsservice.AnalyticsService
}

// NewAnalyticsMiddleware creates an AnalyticsMiddleware.
func NewAnalyticsMiddleware(linkRPC analyticsservice.AnalyticsService) *AnalyticsMiddleware {
	return &AnalyticsMiddleware{linkRPC: linkRPC}
}

// Handle wraps the redirect handler and records visit events asynchronously.
func (m *AnalyticsMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(rw, r)

		if rw.status != http.StatusFound {
			return
		}

		code := strings.TrimPrefix(r.URL.Path, "/")
		ip := httputil.ExtractClientIP(r)
		ua := r.Header.Get("User-Agent")
		ref := r.Header.Get("Referer")
		ts := time.Now().UTC().Format(time.RFC3339Nano)
		rpc := m.linkRPC
		log := logx.WithContext(r.Context())

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			if _, err := rpc.RecordVisit(ctx, &analyticsservice.RecordVisitRequest{
				Code:      code,
				VisitedAt: ts,
				Ip:        ip,
				UserAgent: ua,
				Referer:   ref,
			}); err != nil {
				log.Errorf("record visit async: %v", err)
			}
		}()
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
