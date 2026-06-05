// Package auth contains management authentication helpers for link-api.
package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	// ErrInvalidToken reports an invalid, expired, or malformed management token.
	ErrInvalidToken = errors.New("invalid token")
)

// Config contains JWT token settings for management authentication.
type Config struct {
	Secret          string
	TokenTTLSeconds int64
	Now             func() time.Time
}

// AdminSubject is the authenticated administrator identity stored in a token.
type AdminSubject struct {
	ID       int64
	Username string
}

// TokenManager creates and validates management JWTs.
type TokenManager struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

type claims struct {
	AdminID  int64  `json:"admin_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewTokenManager creates a token manager from configuration.
func NewTokenManager(c Config) *TokenManager {
	now := c.Now
	if now == nil {
		now = time.Now
	}

	return &TokenManager{
		secret: []byte(c.Secret),
		ttl:    time.Duration(c.TokenTTLSeconds) * time.Second,
		now:    now,
	}
}

// Create signs a JWT for the given administrator subject.
func (m *TokenManager) Create(subject AdminSubject) (string, time.Time, error) {
	now := m.now().UTC()
	expiresAt := now.Add(m.ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		AdminID:  subject.ID,
		Username: subject.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	})

	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign management token: %w", err)
	}

	return signed, expiresAt, nil
}

// Validate verifies a JWT and returns its administrator subject.
func (m *TokenManager) Validate(rawToken string) (AdminSubject, error) {
	if rawToken == "" {
		return AdminSubject{}, ErrInvalidToken
	}

	parsedClaims := &claims{}
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, err := parser.ParseWithClaims(rawToken, parsedClaims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return AdminSubject{}, ErrInvalidToken
	}

	if parsedClaims.ExpiresAt == nil || !parsedClaims.ExpiresAt.After(m.now().UTC()) {
		return AdminSubject{}, ErrInvalidToken
	}
	if parsedClaims.AdminID <= 0 || parsedClaims.Username == "" {
		return AdminSubject{}, ErrInvalidToken
	}

	return AdminSubject{
		ID:       parsedClaims.AdminID,
		Username: parsedClaims.Username,
	}, nil
}

