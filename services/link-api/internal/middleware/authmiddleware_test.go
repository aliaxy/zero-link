package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aliaxy/zero-link/services/link-api/internal/auth"
)

func TestAuthMiddleware_RejectsMissingBearerToken(t *testing.T) {
	middleware := NewAuthMiddleware(auth.NewTokenManager(auth.Config{
		Secret:          "local-test-secret",
		TokenTTLSeconds: 3600,
	}))

	handler := middleware.Handle(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("next handler should not be called")
	})

	recorder := httptest.NewRecorder()
	handler(recorder, httptest.NewRequest(http.MethodGet, "/admin/profile", nil))

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestAuthMiddleware_RejectsInvalidBearerToken(t *testing.T) {
	middleware := NewAuthMiddleware(auth.NewTokenManager(auth.Config{
		Secret:          "local-test-secret",
		TokenTTLSeconds: 3600,
	}))

	handler := middleware.Handle(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("next handler should not be called")
	})

	request := httptest.NewRequest(http.MethodGet, "/admin/profile", nil)
	request.Header.Set("Authorization", "Bearer invalid-token")

	recorder := httptest.NewRecorder()
	handler(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestAuthMiddleware_AllowsValidBearerToken(t *testing.T) {
	manager := auth.NewTokenManager(auth.Config{
		Secret:          "local-test-secret",
		TokenTTLSeconds: 3600,
	})
	token, _, err := manager.Create(auth.AdminSubject{
		ID:       42,
		Username: "admin",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	middleware := NewAuthMiddleware(manager)
	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		subject, ok := AdminSubjectFromContext(r.Context())
		if !ok {
			t.Fatal("AdminSubjectFromContext() ok = false, want true")
		}
		if subject.ID != 42 {
			t.Fatalf("subject ID = %d, want 42", subject.ID)
		}
		if subject.Username != "admin" {
			t.Fatalf("subject username = %q, want admin", subject.Username)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/admin/profile", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	handler(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
}
