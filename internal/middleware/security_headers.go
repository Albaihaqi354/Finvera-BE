package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds standard HTTP security headers to every response.
func SecurityHeaders() gin.HandlerFunc {
	isProd := os.Getenv("APP_ENV") == "production"
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		// Prevent MIME-type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		// Enable XSS filter in older browsers
		c.Header("X-XSS-Protection", "1; mode=block")
		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// Permissions policy — disable unused browser features
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		// Content Security Policy — adjust as needed for Swagger UI
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
		// HSTS — only set in production to avoid breaking HTTP-only development
		if isProd {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	}
}
