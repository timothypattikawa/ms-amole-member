package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/argon2"
)

// TokenConfig holds PASETO token configuration
type TokenConfig struct {
	SymmetricKey []byte
	Issuer       string
	AccessTTL    time.Duration
	RefreshTTL   time.Duration
}

// Claims represents the token payload structure
type Claims struct {
	UserID    string    `json:"user_id"`
	Scope     string    `json:"scope"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}

// DeriveKey uses Argon2 to derive a key from a password and salt
func DeriveKey(password, salt []byte) []byte {
	return argon2.IDKey(
		password,
		salt,
		1,       // time cost
		64*1024, // memory in KiB
		4,       // threads
		32,      // key length
	)
}

// GenerateAccessToken creates a new PASETO v2 local (symmetric) token
func (tc *TokenConfig) GenerateAccessToken(userID, scope string) (string, error) {
	// Validate key
	if len(tc.SymmetricKey) < 32 {
		return "", errors.New("symmetric key is too short")
	}

	// Create claims
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		Scope:     scope,
		IssuedAt:  now,
		ExpiresAt: now.Add(tc.AccessTTL),
	}

	// Generate PASETO token
	v2 := paseto.NewV2()
	token, err := v2.Encrypt(tc.SymmetricKey, claims, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

func (tc *TokenConfig) ValidateAccessToken(tokenString string) (*Claims, error) {
	v2 := paseto.NewV2()
	var claims Claims

	// Decrypt and parse the token
	err := v2.Decrypt(tokenString, tc.SymmetricKey, &claims, nil)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %v", err)
	}

	// Additional validation checks
	now := time.Now()
	if now.After(claims.ExpiresAt) {
		return nil, errors.New("token has expired")
	}

	return &claims, nil
}

func (tc *TokenConfig) GenerateRefreshToken(userID string) (string, error) {
	// Create claims with longer expiration
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		Scope:     "refresh",
		IssuedAt:  now,
		ExpiresAt: now.Add(tc.RefreshTTL),
	}

	// Generate PASETO token
	v2 := paseto.NewV2()
	token, err := v2.Encrypt(tc.SymmetricKey, claims, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return token, nil
}
