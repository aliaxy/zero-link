// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"context"
	"net/http"

	"github.com/aliaxy/zero-link/services/link-api/pkg/httputil"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// rateLimiter is the interface satisfied by limit.PeriodLimit.
type rateLimiter interface {
	TakeCtx(ctx context.Context, key string) (int, error)
}

// IPRateLimitMiddleware limits requests per client IP using a Redis sliding window.
type IPRateLimitMiddleware struct {
	limiter        rateLimiter
	trustedProxies []string
}

// NewIPRateLimitMiddleware creates an IP rate limiter with the given window period (seconds),
// quota, Redis store, key prefix, and trusted proxy IPs.
func NewIPRateLimitMiddleware(store *redis.Redis, period, quota int, keyPrefix string, trustedProxies []string) *IPRateLimitMiddleware {
	return &IPRateLimitMiddleware{
		limiter:        limit.NewPeriodLimit(period, quota, store, keyPrefix),
		trustedProxies: trustedProxies,
	}
}

// Handle rejects requests that exceed the per-IP quota with 429 Too Many Requests.
func (m *IPRateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, err := m.limiter.TakeCtx(r.Context(), httputil.ExtractClientIP(r, m.trustedProxies))
		if err != nil || code == limit.OverQuota {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}
