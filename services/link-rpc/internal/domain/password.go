package domain

import "golang.org/x/crypto/bcrypt"

// VerifyPassword checks a plaintext password against a bcrypt hash.
func VerifyPassword(hash, password string) bool {
	if hash == "" || password == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
