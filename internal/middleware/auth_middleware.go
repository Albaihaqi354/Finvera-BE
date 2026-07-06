package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"finvera-be/internal/dto"
	"finvera-be/pkg/blacklist"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates the Bearer JWT token on every protected route.
// It sets "userId" (string UUID) in the Gin context for downstream handlers.
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Authorization header is required"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Authorization header format must be: Bearer {token}"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		if blacklist.IsBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Token has been revoked"))
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Reject any token not using HMAC — prevents alg:none attacks
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid or expired token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Failed to parse token claims"))
			c.Abort()
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Token does not contain a valid user ID"))
			c.Abort()
			return
		}

		c.Set("userId", userID)
		c.Next()
	}
}
