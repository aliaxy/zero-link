// Package domain contains Stage 3 management domain rules.
package domain

import (
	"crypto/rand"
	"errors"
	"math/big"
	"net/url"
	"regexp"
	"time"
)

const (
	// GeneratedCodeLength is the Stage 3 automatic short-code length.
	GeneratedCodeLength = 6
	// DefaultPage is the default management list page.
	DefaultPage = 1
	// DefaultPageSize is the default management list page size.
	DefaultPageSize = 20
	// MaxPageSize is the maximum management list page size.
	MaxPageSize = 100
	// LinkStatusActive marks a short link active for management.
	LinkStatusActive = 1
	// LinkStatusDisabled marks a short link disabled for management.
	LinkStatusDisabled = 2
)

var (
	// ErrInvalidCode reports an invalid custom short code.
	ErrInvalidCode = errors.New("invalid short code")
	// ErrReservedCode reports a custom short code that conflicts with system routes.
	ErrReservedCode = errors.New("reserved short code")
	// ErrInvalidOriginURL reports an invalid destination URL.
	ErrInvalidOriginURL = errors.New("invalid origin url")
	// ErrInvalidExpiration reports an invalid expiration timestamp.
	ErrInvalidExpiration = errors.New("invalid expiration")
	// ErrInvalidPagination reports invalid list pagination parameters.
	ErrInvalidPagination = errors.New("invalid pagination")
	// ErrInvalidStatus reports an invalid short-link status.
	ErrInvalidStatus = errors.New("invalid status")

	customCodePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{3,32}$`)
	reservedCodes     = map[string]struct{}{
		"admin":   {},
		"api":     {},
		"healthz": {},
		"readyz":  {},
		"rpc":     {},
		"metrics": {},
		"static":  {},
	}
	base62Alphabet = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
)

// GenerateCode returns a cryptographically random base62 short code.
func GenerateCode() (string, error) {
	buf := make([]byte, GeneratedCodeLength)
	alphabetSize := big.NewInt(int64(len(base62Alphabet)))

	for i := range buf {
		n, err := rand.Int(rand.Reader, alphabetSize)
		if err != nil {
			return "", err
		}
		buf[i] = base62Alphabet[n.Int64()]
	}

	return string(buf), nil
}

// ValidateCustomCode checks Stage 3 custom short-code syntax and reserved words.
func ValidateCustomCode(code string) error {
	if !customCodePattern.MatchString(code) {
		return ErrInvalidCode
	}
	if _, ok := reservedCodes[code]; ok {
		return ErrReservedCode
	}
	return nil
}

// ValidateOriginURL checks that the destination URL is absolute HTTP or HTTPS.
func ValidateOriginURL(rawURL string) error {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return ErrInvalidOriginURL
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrInvalidOriginURL
	}
	if parsed.Host == "" {
		return ErrInvalidOriginURL
	}
	return nil
}

// ValidateExpireAt parses and validates an optional future expiration timestamp.
func ValidateExpireAt(rawExpireAt string, now time.Time) (time.Time, error) {
	if rawExpireAt == "" {
		return time.Time{}, nil
	}

	expireAt, err := time.Parse(time.RFC3339, rawExpireAt)
	if err != nil {
		return time.Time{}, ErrInvalidExpiration
	}
	if !expireAt.After(now) {
		return time.Time{}, ErrInvalidExpiration
	}

	return expireAt, nil
}

// NormalizePagination applies defaults and validates management pagination.
func NormalizePagination(page, pageSize int64) (int64, int64, error) {
	if page == 0 {
		page = DefaultPage
	}
	if pageSize == 0 {
		pageSize = DefaultPageSize
	}
	if page < 1 || pageSize < 1 || pageSize > MaxPageSize {
		return 0, 0, ErrInvalidPagination
	}
	return page, pageSize, nil
}

// ValidateLinkStatus checks Stage 3 short-link management status values.
func ValidateLinkStatus(status int64) error {
	if status == LinkStatusActive || status == LinkStatusDisabled {
		return nil
	}
	return ErrInvalidStatus
}
