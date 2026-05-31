// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// IPRateLimitMiddleware limits requests per client IP using a Redis sliding window.
type IPRateLimitMiddleware struct {
	limiter *limit.PeriodLimit
}

// NewIPRateLimitMiddleware creates an IP rate limiter with the given window period (seconds),
// quota, Redis store, and key prefix.
func NewIPRateLimitMiddleware(store *redis.Redis, period, quota int, keyPrefix string) *IPRateLimitMiddleware {
	return &IPRateLimitMiddleware{
		limiter: limit.NewPeriodLimit(period, quota, store, keyPrefix),
	}
}

// Handle rejects requests that exceed the per-IP quota with 429 Too Many Requests.
func (m *IPRateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := extractClientIP(r)
		code, err := m.limiter.TakeCtx(r.Context(), ip)
		if err != nil || code == limit.OverQuota {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

func extractClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
