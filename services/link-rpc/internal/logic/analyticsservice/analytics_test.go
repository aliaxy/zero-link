package analyticsservicelogic

import (
	"testing"
	"time"
)

func TestHashIP_Deterministic(t *testing.T) {
	h1 := hashIP("1.2.3.4", "salt")
	h2 := hashIP("1.2.3.4", "salt")
	if h1 != h2 {
		t.Fatalf("hashIP not deterministic: %q != %q", h1, h2)
	}
}

func TestHashIP_DifferentSalt(t *testing.T) {
	h1 := hashIP("1.2.3.4", "salt-a")
	h2 := hashIP("1.2.3.4", "salt-b")
	if h1 == h2 {
		t.Fatal("hashIP with different salts should produce different hashes")
	}
}

func TestDetectDevice(t *testing.T) {
	tests := []struct {
		name string
		ua   string
		want string
	}{
		{"empty ua", "", "unknown"},
		{"googlebot", "Googlebot/2.1 (+http://www.google.com/bot.html)", "bot"},
		{"crawler", "Mozilla/5.0 (compatible; bingbot/2.0)", "bot"},
		{"iphone", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X)", "mobile"},
		{"android", "Mozilla/5.0 (Linux; Android 13; Pixel 7)", "mobile"},
		{"mobile keyword", "Mozilla/5.0 (Linux; Android 13) Mobile Safari/537.36", "mobile"},
		{"desktop chrome", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0.0.0", "desktop"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectDevice(tt.ua)
			if got != tt.want {
				t.Fatalf("detectDevice(%q) = %q, want %q", tt.ua, got, tt.want)
			}
		})
	}
}

func TestParseDateRange_Defaults(t *testing.T) {
	from, to, err := parseDateRange("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	now := time.Now().UTC().Truncate(24 * time.Hour)
	if !to.Equal(now) {
		t.Fatalf("default to = %v, want %v", to, now)
	}
	wantFrom := now.AddDate(0, 0, -defaultRangeDays)
	if !from.Equal(wantFrom) {
		t.Fatalf("default from = %v, want %v", from, wantFrom)
	}
}

func TestParseDateRange_InvalidFrom(t *testing.T) {
	_, _, err := parseDateRange("not-a-date", "")
	if err == nil {
		t.Fatal("expected error for invalid from date")
	}
}

func TestParseDateRange_InvalidTo(t *testing.T) {
	_, _, err := parseDateRange("", "not-a-date")
	if err == nil {
		t.Fatal("expected error for invalid to date")
	}
}

func TestParseDateRange_FromAfterTo(t *testing.T) {
	_, _, err := parseDateRange("2026-05-10", "2026-05-01")
	if err == nil {
		t.Fatal("expected error when from is after to")
	}
}

func TestParseDateRange_ExceedsMax(t *testing.T) {
	_, _, err := parseDateRange("2026-01-01", "2026-05-01")
	if err == nil {
		t.Fatal("expected error when range exceeds 90 days")
	}
}

func TestParseDateRange_Valid(t *testing.T) {
	from, to, err := parseDateRange("2026-05-01", "2026-05-29")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if from.Format("2006-01-02") != "2026-05-01" {
		t.Fatalf("from = %v, want 2026-05-01", from)
	}
	if to.Format("2006-01-02") != "2026-05-29" {
		t.Fatalf("to = %v, want 2026-05-29", to)
	}
}
