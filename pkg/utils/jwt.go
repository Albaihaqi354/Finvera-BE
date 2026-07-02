package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Token durations
const (
	AccessTokenDuration = 2 * time.Hour // short-lived access token
)

type JWTClaim struct {
	UserID   uuid.UUID `json:"userId"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a signed HMAC-SHA256 JWT access token.
func GenerateJWT(userID uuid.UUID, username, secret, issuer string) (string, error) {
	now := time.Now()
	claims := &JWTClaim{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
