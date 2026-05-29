package middleware

import (
	"net/http"
	"strings"
)

// CorsMiddleware adds CORS headers for configured origins.
type CorsMiddleware struct {
	allowOrigins map[string]struct{}
}

// NewCorsMiddleware creates a CorsMiddleware from a list of allowed origin strings.
func NewCorsMiddleware(origins []string) *CorsMiddleware {
	m := make(map[string]struct{}, len(origins))
	for _, o := range origins {
		m[strings.TrimRight(o, "/")] = struct{}{}
	}
	return &CorsMiddleware{allowOrigins: m}
}

// Handle implements go-zero middleware. Applied globally via server.Use() so that
// OPTIONS preflight requests are handled before route matching.
func (m *CorsMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			if _, ok := m.allowOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
				w.Header().Set("Vary", "Origin")
			}
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}
