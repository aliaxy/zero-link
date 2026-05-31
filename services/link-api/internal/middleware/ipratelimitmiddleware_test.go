package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zeromicro/go-zero/core/limit"
)

// stubLimiter is a fake rateLimiter for unit tests.
type stubLimiter struct {
	code int
	err  error
}

func (s *stubLimiter) TakeCtx(_ context.Context, _ string) (int, error) {
	return s.code, s.err
}

func newIPMiddlewareWithStub(code int, err error) *IPRateLimitMiddleware {
	return &IPRateLimitMiddleware{limiter: &stubLimiter{code: code, err: err}}
}

func TestIPRateLimitMiddleware_AllowsRequest(t *testing.T) {
	mw := newIPMiddlewareWithStub(limit.Allowed, nil)

	called := false
	handler := mw.Handle(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	handler(recorder, httptest.NewRequest(http.MethodGet, "/abc", nil))

	if !called {
		t.Fatal("next handler should be called when Allowed")
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
}

func TestIPRateLimitMiddleware_BlocksOnOverQuota(t *testing.T) {
	mw := newIPMiddlewareWithStub(limit.OverQuota, nil)

	handler := mw.Handle(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("next handler should not be called when OverQuota")
	})

	recorder := httptest.NewRecorder()
	handler(recorder, httptest.NewRequest(http.MethodGet, "/abc", nil))

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusTooManyRequests)
	}
}

func TestIPRateLimitMiddleware_BlocksOnLimiterError(t *testing.T) {
	mw := newIPMiddlewareWithStub(limit.Unknown, errors.New("redis unavailable"))

	handler := mw.Handle(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("next handler should not be called on limiter error")
	})

	recorder := httptest.NewRecorder()
	handler(recorder, httptest.NewRequest(http.MethodGet, "/abc", nil))

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusTooManyRequests)
	}
}

func TestIPRateLimitMiddleware_AllowsOnHitQuota(t *testing.T) {
	// HitQuota means this request exactly hit the quota — still allowed.
	mw := newIPMiddlewareWithStub(limit.HitQuota, nil)

	called := false
	handler := mw.Handle(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	handler(recorder, httptest.NewRequest(http.MethodGet, "/abc", nil))

	if !called {
		t.Fatal("next handler should be called when HitQuota")
	}
}

func TestIPRateLimitMiddleware_UsesXForwardedForAsKey(t *testing.T) {
	var capturedKey string
	mw := &IPRateLimitMiddleware{
		limiter: &capturingLimiter{capture: &capturedKey},
	}

	req := httptest.NewRequest(http.MethodGet, "/abc", nil)
	req.Header.Set("X-Forwarded-For", "1.2.3.4")

	mw.Handle(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})(httptest.NewRecorder(), req)

	if capturedKey != "1.2.3.4" {
		t.Fatalf("limiter key = %q, want 1.2.3.4", capturedKey)
	}
}

func TestIPRateLimitMiddleware_FallsBackToRemoteAddr(t *testing.T) {
	var capturedKey string
	mw := &IPRateLimitMiddleware{
		limiter: &capturingLimiter{capture: &capturedKey},
	}

	req := httptest.NewRequest(http.MethodGet, "/abc", nil)
	req.RemoteAddr = "5.6.7.8:1234"

	mw.Handle(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})(httptest.NewRecorder(), req)

	if capturedKey != "5.6.7.8" {
		t.Fatalf("limiter key = %q, want 5.6.7.8", capturedKey)
	}
}

// capturingLimiter records the key passed to TakeCtx and always allows.
type capturingLimiter struct {
	capture *string
}

func (c *capturingLimiter) TakeCtx(_ context.Context, key string) (int, error) {
	*c.capture = key
	return limit.Allowed, nil
}
