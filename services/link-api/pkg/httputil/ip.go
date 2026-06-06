// Package httputil provides shared HTTP helper utilities for link-api.
package httputil

import (
	"net/http"
	"slices"
	"strings"
)

// ExtractClientIP returns the real client IP from the request.
// X-Forwarded-For is only trusted when the direct RemoteAddr is in trustedProxies.
// Pass nil or an empty slice to always use RemoteAddr directly.
func ExtractClientIP(r *http.Request, trustedProxies []string) string {
	direct := directIP(r.RemoteAddr)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" && isTrustedProxy(direct, trustedProxies) {
		if ip, _, ok := strings.Cut(xff, ","); ok {
			return strings.TrimSpace(ip)
		}
		return strings.TrimSpace(xff)
	}
	return direct
}

func directIP(remoteAddr string) string {
	if idx := strings.LastIndex(remoteAddr, ":"); idx != -1 {
		return remoteAddr[:idx]
	}
	return remoteAddr
}

func isTrustedProxy(ip string, proxies []string) bool {
	return slices.Contains(proxies, ip)
}
