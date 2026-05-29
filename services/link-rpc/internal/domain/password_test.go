package domain

import "testing"

func TestVerifyPassword(t *testing.T) {
	const seededHash = "$2a$10$UpD2JjqWVgQOatvqxd5H3OSQwzxC5o5gYf31R73AJIz.dQOAuKkBS"

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{name: "seeded password", password: "zerolink", want: true},
		{name: "wrong password", password: "wrong-password"},
		{name: "empty password"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VerifyPassword(seededHash, tt.password)
			if got != tt.want {
				t.Fatalf("VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyPasswordRejectsInvalidHash(t *testing.T) {
	if VerifyPassword("not-a-bcrypt-hash", "zerolink") {
		t.Fatal("VerifyPassword() = true, want false")
	}
}
