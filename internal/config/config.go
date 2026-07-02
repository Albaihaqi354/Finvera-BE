package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBDSN          string
	JWTSecret      string
	JWTIssuer      string
	AppEnv         string
	AllowedOrigins []string
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func LoadConfig() *Config {
	// Try to load .env from current dir and up to 3 parent dirs.
	// This handles running from cmd/server/ or from project root.
	loaded := false
	dir, _ := os.Getwd()
	for i := 0; i < 4; i++ {
		envPath := filepath.Join(dir, ".env")
		if err := godotenv.Load(envPath); err == nil {
			log.Printf("Loaded .env from: %s", envPath)
			loaded = true
			break
		}
		dir = filepath.Dir(dir)
	}
	if !loaded {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" {
		log.Fatal("Database configuration is incomplete. Check DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME in .env")
	}

	// Determine SSL mode based on environment
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}
	sslMode := "disable"
	if appEnv == "production" {
		sslMode = "require"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		dbHost, dbUser, dbPass, dbName, dbPort, sslMode,
	)

	// Validate JWT secret — must be at least 32 chars for HMAC-SHA256 security
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set. Generate with: openssl rand -hex 64")
	}
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET is too short. Use at least 32 characters (64 hex chars recommended)")
	}

	// Parse allowed origins
	allowedOriginsRaw := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsRaw != "" {
		for _, o := range strings.Split(allowedOriginsRaw, ",") {
			trimmed := strings.TrimSpace(o)
			if trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	}
	if len(allowedOrigins) == 0 {
		if appEnv == "production" {
			log.Fatal("ALLOWED_ORIGINS must be set in production environment")
		}
		// Default for development only
		allowedOrigins = []string{"http://localhost:3000"}
		log.Println("Warning: ALLOWED_ORIGINS not set. Defaulting to http://localhost:3000 (development only)")
	}

	jwtIssuer := os.Getenv("JWT_ISSUER")
	if jwtIssuer == "" {
		jwtIssuer = "finvera"
	}

	return &Config{
		Port:           port,
		DBDSN:          dsn,
		JWTSecret:      jwtSecret,
		JWTIssuer:      jwtIssuer,
		AppEnv:         appEnv,
		AllowedOrigins: allowedOrigins,
	}
}
