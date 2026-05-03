package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
	secret string
}

func NewPasswordService(secret string) *PasswordService {
	return &PasswordService{secret: secret}
}

// prehash combines password with secret and returns a fixed-length SHA-256 hash
// This ensures we never exceed bcrypt's 72-byte limit
func (p *PasswordService) prehash(password string) string {
	hash := sha256.Sum256([]byte(password + p.secret))
	return hex.EncodeToString(hash[:])
}

// HashPassword hashes a password with the secret and bcrypt
func (p *PasswordService) HashPassword(password string) (string, error) {
	// Pre-hash to get fixed 64-byte hex string (within bcrypt's 72-byte limit)
	prehashed := p.prehash(password)
	bytes, err := bcrypt.GenerateFromPassword([]byte(prehashed), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a password with a hash
func (p *PasswordService) CheckPassword(password, hash string) bool {
	prehashed := p.prehash(password)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(prehashed))
	return err == nil
}
