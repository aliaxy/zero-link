// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// LoginRateLimitMiddleware limits login attempts per client IP to prevent brute force attacks.
// It delegates to IPRateLimitMiddleware with a per-minute window.
type LoginRateLimitMiddleware struct {
	inner *IPRateLimitMiddleware
}

// NewLoginRateLimitMiddleware creates a login rate limiter capped at quota attempts per minute per IP.
func NewLoginRateLimitMiddleware(store *redis.Redis, quota int) *LoginRateLimitMiddleware {
	return &LoginRateLimitMiddleware{
		inner: NewIPRateLimitMiddleware(store, 60, quota, "rl:login:ip:"),
	}
}

// Handle rejects excessive login attempts with 429 Too Many Requests.
func (m *LoginRateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return m.inner.Handle(next)
}
