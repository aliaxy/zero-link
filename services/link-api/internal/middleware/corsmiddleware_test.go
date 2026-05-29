package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newCorsHandler(origins []string, next http.HandlerFunc) http.HandlerFunc {
	return NewCorsMiddleware(origins).Handle(next)
}

func noop(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }

func TestCorsMiddleware_AllowedOrigin_SetsHeaders(t *testing.T) {
	handler := newCorsHandler([]string{"http://localhost:5173"}, noop)

	r := httptest.NewRequest(http.MethodGet, "/admin/links", nil)
	r.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()
	handler(w, r)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("Allow-Origin = %q, want %q", got, "http://localhost:5173")
	}
	if got := w.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Fatal("Allow-Methods header missing")
	}
	if got := w.Header().Get("Access-Control-Allow-Headers"); got == "" {
		t.Fatal("Allow-Headers header missing")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCorsMiddleware_DisallowedOrigin_NoHeaders(t *testing.T) {
	handler := newCorsHandler([]string{"http://localhost:5173"}, noop)

	r := httptest.NewRequest(http.MethodGet, "/admin/links", nil)
	r.Header.Set("Origin", "http://evil.example.com")
	w := httptest.NewRecorder()
	handler(w, r)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("Allow-Origin = %q, want empty", got)
	}
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCorsMiddleware_NoOrigin_PassesThrough(t *testing.T) {
	called := false
	handler := newCorsHandler([]string{"http://localhost:5173"}, func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	handler(w, r)

	if !called {
		t.Fatal("next handler should be called when no Origin header")
	}
}

func TestCorsMiddleware_OptionsPreflight_Returns204(t *testing.T) {
	handler := newCorsHandler([]string{"http://localhost:5173"}, func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("next handler should not be called for OPTIONS preflight")
	})

	r := httptest.NewRequest(http.MethodOptions, "/admin/login", nil)
	r.Header.Set("Origin", "http://localhost:5173")
	r.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()
	handler(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("Allow-Origin = %q, want %q", got, "http://localhost:5173")
	}
}

func TestCorsMiddleware_OptionsDisallowedOrigin_NoHeadersAnd204(t *testing.T) {
	handler := newCorsHandler([]string{"http://localhost:5173"}, func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("next handler should not be called for OPTIONS")
	})

	r := httptest.NewRequest(http.MethodOptions, "/admin/login", nil)
	r.Header.Set("Origin", "http://evil.example.com")
	w := httptest.NewRecorder()
	handler(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("Allow-Origin = %q, want empty", got)
	}
}

func TestCorsMiddleware_MultipleOrigins(t *testing.T) {
	origins := []string{"http://localhost:5173", "https://admin.example.com"}
	handler := newCorsHandler(origins, noop)

	for _, origin := range origins {
		r := httptest.NewRequest(http.MethodGet, "/admin/links", nil)
		r.Header.Set("Origin", origin)
		w := httptest.NewRecorder()
		handler(w, r)

		if got := w.Header().Get("Access-Control-Allow-Origin"); got != origin {
			t.Fatalf("origin %q: Allow-Origin = %q, want %q", origin, got, origin)
		}
	}
}
