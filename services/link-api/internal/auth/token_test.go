package auth

import (
	"testing"
	"time"
)

func TestTokenManager_CreateAndValidate(t *testing.T) {
	now := time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC)
	manager := NewTokenManager(Config{
		Secret:          "local-test-secret",
		TokenTTLSeconds: 3600,
		Now:             func() time.Time { return now },
	})

	token, expiresAt, err := manager.Create(AdminSubject{
		ID:       42,
		Username: "admin",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if want := now.Add(time.Hour); !expiresAt.Equal(want) {
		t.Fatalf("expiresAt = %v, want %v", expiresAt, want)
	}

	got, err := manager.Validate(token)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if got.ID != 42 {
		t.Fatalf("subject ID = %d, want 42", got.ID)
	}
	if got.Username != "admin" {
		t.Fatalf("subject username = %q, want admin", got.Username)
	}
}

func TestTokenManager_ValidateRejectsExpiredToken(t *testing.T) {
	issuedAt := time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC)
	manager := NewTokenManager(Config{
		Secret:          "local-test-secret",
		TokenTTLSeconds: 1,
		Now:             func() time.Time { return issuedAt },
	})

	token, _, err := manager.Create(AdminSubject{
		ID:       42,
		Username: "admin",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	expiredManager := NewTokenManager(Config{
		Secret:          "local-test-secret",
		TokenTTLSeconds: 1,
		Now:             func() time.Time { return issuedAt.Add(2 * time.Second) },
	})

	if _, err := expiredManager.Validate(token); err == nil {
		t.Fatal("Validate() error = nil, want expired token error")
	}
}

func TestTokenManager_ValidateRejectsWrongSecret(t *testing.T) {
	manager := NewTokenManager(Config{
		Secret:          "local-test-secret",
		TokenTTLSeconds: 3600,
	})

	token, _, err := manager.Create(AdminSubject{
		ID:       42,
		Username: "admin",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	wrongManager := NewTokenManager(Config{
		Secret:          "different-secret",
		TokenTTLSeconds: 3600,
	})

	if _, err := wrongManager.Validate(token); err == nil {
		t.Fatal("Validate() error = nil, want signature error")
	}
}

func TestTokenManager_CreateRejectsInvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "missing secret",
			config: Config{
				TokenTTLSeconds: 3600,
			},
		},
		{
			name: "missing ttl",
			config: Config{
				Secret: "local-test-secret",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewTokenManager(tt.config)
			_, _, err := manager.Create(AdminSubject{
				ID:       42,
				Username: "admin",
			})
			if err == nil {
				t.Fatal("Create() error = nil, want config error")
			}
		})
	}
}
