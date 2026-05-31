// Package httputil provides shared HTTP helper utilities for link-api.
package httputil

import (
	"net/http"
	"strings"
)

// ExtractClientIP returns the real client IP from the request.
// It prefers the first value in X-Forwarded-For and falls back to RemoteAddr.
func ExtractClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip, _, ok := strings.Cut(xff, ","); ok {
			return strings.TrimSpace(ip)
		}
		return strings.TrimSpace(xff)
	}
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}
