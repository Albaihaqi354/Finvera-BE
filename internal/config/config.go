package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBDSN     string
	JWTSecret string
	JWTIssuer string
}

func LoadConfig() *Config {
	// Try to load .env from current dir and up to 3 parent dirs
	// This handles running from cmd/server/ or from project root
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

	// Validate required fields
	if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" {
		log.Fatal("Database configuration is incomplete. Check DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME in .env")
	}

	// Construct DSN in proper postgres format
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, dbUser, dbPass, dbName, dbPort)

	return &Config{
		Port:      port,
		DBDSN:     dsn,
		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTIssuer: os.Getenv("JWT_ISSUER"),
	}
}
