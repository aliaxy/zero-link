package domain

import (
	"strings"
	"testing"
	"time"
)

func TestValidateCustomCode(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{name: "minimum length", code: "abc"},
		{name: "maximum length", code: strings.Repeat("a", 32)},
		{name: "letters digits underscore hyphen", code: "Campaign_2026-01"},
		{name: "too short", code: "ab", wantErr: true},
		{name: "too long", code: strings.Repeat("a", 33), wantErr: true},
		{name: "space", code: "bad code", wantErr: true},
		{name: "slash", code: "bad/code", wantErr: true},
		{name: "non ascii", code: "中文", wantErr: true},
		{name: "reserved admin", code: "admin", wantErr: true},
		{name: "reserved healthz", code: "healthz", wantErr: true},
		{name: "reserved metrics", code: "metrics", wantErr: true},
		{name: "reserved static", code: "static", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCustomCode(tt.code)
			if tt.wantErr && err == nil {
				t.Fatal("ValidateCustomCode() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("ValidateCustomCode() error = %v, want nil", err)
			}
		})
	}
}

func TestGenerateCode(t *testing.T) {
	code, err := GenerateCode()
	if err != nil {
		t.Fatalf("GenerateCode() error = %v", err)
	}
	if len(code) != 6 {
		t.Fatalf("len(code) = %d, want 6", len(code))
	}
	if err := ValidateCustomCode(code); err != nil {
		t.Fatalf("generated code is invalid: %v", err)
	}
}

func TestValidateOriginURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantErr bool
	}{
		{name: "https", rawURL: "https://example.com/page"},
		{name: "http", rawURL: "http://example.com/page"},
		{name: "empty", rawURL: "", wantErr: true},
		{name: "relative", rawURL: "/page", wantErr: true},
		{name: "unsupported scheme", rawURL: "ftp://example.com/page", wantErr: true},
		{name: "missing host", rawURL: "https:///page", wantErr: true},
		{name: "malformed", rawURL: "://bad", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOriginURL(tt.rawURL)
			if tt.wantErr && err == nil {
				t.Fatal("ValidateOriginURL() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("ValidateOriginURL() error = %v, want nil", err)
			}
		})
	}
}

func TestValidateExpireAt(t *testing.T) {
	now := time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		expireAt string
		wantErr  bool
	}{
		{name: "empty"},
		{name: "future", expireAt: "2026-12-31T23:59:59Z"},
		{name: "past", expireAt: "2026-01-01T00:00:00Z", wantErr: true},
		{name: "invalid", expireAt: "tomorrow", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateExpireAt(tt.expireAt, now)
			if tt.wantErr && err == nil {
				t.Fatal("ValidateExpireAt() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("ValidateExpireAt() error = %v, want nil", err)
			}
		})
	}
}

func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name         string
		page         int64
		pageSize     int64
		wantPage     int64
		wantPageSize int64
		wantErr      bool
	}{
		{name: "defaults", wantPage: 1, wantPageSize: 20},
		{name: "explicit", page: 2, pageSize: 50, wantPage: 2, wantPageSize: 50},
		{name: "max page size", page: 1, pageSize: 100, wantPage: 1, wantPageSize: 100},
		{name: "invalid page", page: -1, pageSize: 20, wantErr: true},
		{name: "invalid page size", page: 1, pageSize: 101, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotPageSize, err := NormalizePagination(tt.page, tt.pageSize)
			if tt.wantErr && err == nil {
				t.Fatal("NormalizePagination() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("NormalizePagination() error = %v, want nil", err)
			}
			if tt.wantErr {
				return
			}
			if gotPage != tt.wantPage {
				t.Fatalf("page = %d, want %d", gotPage, tt.wantPage)
			}
			if gotPageSize != tt.wantPageSize {
				t.Fatalf("pageSize = %d, want %d", gotPageSize, tt.wantPageSize)
			}
		})
	}
}
