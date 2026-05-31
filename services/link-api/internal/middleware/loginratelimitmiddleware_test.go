package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zeromicro/go-zero/core/limit"
)

func newLoginMiddlewareWithStub(code int) *LoginRateLimitMiddleware {
	return &LoginRateLimitMiddleware{
		inner: &IPRateLimitMiddleware{limiter: &stubLimiter{code: code}},
	}
}

func TestLoginRateLimitMiddleware_AllowsRequest(t *testing.T) {
	mw := newLoginMiddlewareWithStub(limit.Allowed)

	called := false
	handler := mw.Handle(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	handler(recorder, httptest.NewRequest(http.MethodPost, "/admin/login", nil))

	if !called {
		t.Fatal("next handler should be called when Allowed")
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
}

func TestLoginRateLimitMiddleware_BlocksOnOverQuota(t *testing.T) {
	mw := newLoginMiddlewareWithStub(limit.OverQuota)

	handler := mw.Handle(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("next handler should not be called when OverQuota")
	})

	recorder := httptest.NewRecorder()
	handler(recorder, httptest.NewRequest(http.MethodPost, "/admin/login", nil))

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusTooManyRequests)
	}
}
