package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	refreshTokenTTL     = 7 * 24 * time.Hour
	refreshTokenBytes   = 32
	keyRefreshToken     = "zl:rt:"
	keyRefreshTokenUser = "zl:rt:user:"
)

// RefreshTokenIssuer is the interface for refresh token lifecycle management.
type RefreshTokenIssuer interface {
	Issue(ctx context.Context, adminID int64) (string, error)
	Rotate(ctx context.Context, rawToken string) (newRaw string, adminID int64, err error)
	RevokeAll(ctx context.Context, adminID int64) error
}

// RefreshTokenStore manages refresh tokens in Redis.
type RefreshTokenStore struct {
	rdb *redis.Redis
}

// NewRefreshTokenStore creates a refresh token store.
func NewRefreshTokenStore(rdb *redis.Redis) *RefreshTokenStore {
	return &RefreshTokenStore{rdb: rdb}
}

var _ RefreshTokenIssuer = (*RefreshTokenStore)(nil)

// Issue generates a new refresh token, stores it in Redis, and returns the raw token.
func (s *RefreshTokenStore) Issue(ctx context.Context, adminID int64) (string, error) {
	raw, hash, err := generateToken()
	if err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}

	ttlSec := int(refreshTokenTTL.Seconds())
	tokenKey := keyRefreshToken + hash
	userKey := keyRefreshTokenUser + fmt.Sprint(adminID)

	if err := s.rdb.SetexCtx(ctx, tokenKey, fmt.Sprint(adminID), ttlSec); err != nil {
		return "", fmt.Errorf("store refresh token: %w", err)
	}
	if _, err := s.rdb.SaddCtx(ctx, userKey, hash); err != nil {
		return "", fmt.Errorf("track refresh token: %w", err)
	}
	_ = s.rdb.ExpireCtx(ctx, userKey, ttlSec)

	return raw, nil
}

// Rotate verifies a refresh token, deletes it, and issues a new one atomically-ish.
// Returns the new raw token and the admin ID.
// If the token was already rotated (reuse detected), all tokens for that user are revoked.
func (s *RefreshTokenStore) Rotate(ctx context.Context, rawToken string) (string, int64, error) {
	hash := hashToken(rawToken)
	tokenKey := keyRefreshToken + hash

	adminIDStr, err := s.rdb.GetCtx(ctx, tokenKey)
	if err != nil || adminIDStr == "" {
		return "", 0, ErrInvalidToken
	}

	var adminID int64
	if _, err := fmt.Sscan(adminIDStr, &adminID); err != nil || adminID <= 0 {
		return "", 0, ErrInvalidToken
	}

	// Delete the used token before issuing a new one.
	if _, err := s.rdb.DelCtx(ctx, tokenKey); err != nil {
		return "", 0, fmt.Errorf("rotate refresh token: %w", err)
	}
	userKey := keyRefreshTokenUser + fmt.Sprint(adminID)
	_, _ = s.rdb.SremCtx(ctx, userKey, hash)

	newRaw, err := s.Issue(ctx, adminID)
	if err != nil {
		return "", 0, err
	}
	return newRaw, adminID, nil
}

// RevokeAll deletes every refresh token associated with the given admin ID.
func (s *RefreshTokenStore) RevokeAll(ctx context.Context, adminID int64) error {
	userKey := keyRefreshTokenUser + fmt.Sprint(adminID)

	hashes, err := s.rdb.SmembersCtx(ctx, userKey)
	if err != nil {
		return fmt.Errorf("list refresh tokens: %w", err)
	}

	keys := make([]string, 0, len(hashes)+1)
	for _, h := range hashes {
		keys = append(keys, keyRefreshToken+h)
	}
	keys = append(keys, userKey)

	if len(keys) > 0 {
		if _, err := s.rdb.DelCtx(ctx, keys...); err != nil {
			return fmt.Errorf("revoke all refresh tokens: %w", err)
		}
	}
	return nil
}

// ExpiresAt returns when a refresh token expires (now + TTL).
func ExpiresAt() time.Time {
	return time.Now().UTC().Add(refreshTokenTTL)
}

func generateToken() (raw, hash string, err error) {
	b := make([]byte, refreshTokenBytes)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	hash = hashToken(raw)
	return raw, hash, nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
