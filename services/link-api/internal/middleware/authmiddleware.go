// Package middleware contains link-api HTTP middleware.
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aliaxy/zero-link/services/link-api/internal/auth"
	"github.com/zeromicro/go-zero/core/logx"
)

type contextKey string

const adminSubjectContextKey contextKey = "admin-subject"

// AuthMiddleware validates management Bearer tokens.
type AuthMiddleware struct {
	tokenManager *auth.TokenManager
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewAuthMiddleware creates management authentication middleware.
func NewAuthMiddleware(tokenManager *auth.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
	}
}

// Handle rejects unauthenticated requests and stores the admin subject on success.
func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subject, err := m.authenticate(r)
		if err != nil {
			writeUnauthorized(w)
			return
		}

		logx.WithContext(r.Context()).Infow("admin authenticated",
			logx.Field("admin_id", subject.ID),
			logx.Field("username", subject.Username),
		)
		ctx := ContextWithAdminSubject(r.Context(), subject)
		next(w, r.WithContext(ctx))
	}
}

// AdminSubjectFromContext returns the authenticated administrator subject from context.
func AdminSubjectFromContext(ctx context.Context) (auth.AdminSubject, bool) {
	subject, ok := ctx.Value(adminSubjectContextKey).(auth.AdminSubject)
	return subject, ok
}

// ContextWithAdminSubject stores the authenticated administrator subject in context.
func ContextWithAdminSubject(ctx context.Context, subject auth.AdminSubject) context.Context {
	return context.WithValue(ctx, adminSubjectContextKey, subject)
}

func (m *AuthMiddleware) authenticate(r *http.Request) (auth.AdminSubject, error) {
	const prefix = "Bearer "

	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, prefix) {
		return auth.AdminSubject{}, auth.ErrInvalidToken
	}

	return m.tokenManager.Validate(strings.TrimSpace(strings.TrimPrefix(header, prefix)))
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(errorResponse{
		Code:    "UNAUTHENTICATED",
		Message: "missing or invalid bearer token",
	})
}
