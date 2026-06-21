package main

import (
	"finvera-be/internal/config"
	"finvera-be/internal/database"
	"finvera-be/internal/router"
	"log"

	_ "finvera-be/docs" // swagger docs - wajib di-import

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title           Finvera API
// @version         1.0
// @description     REST API untuk aplikasi keuangan Finvera - Personal Finance Manager

// @contact.name   Finvera Dev Team
// @contact.email  dev@finvera.app

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token JWT dengan format: Bearer {token}

func main() {
	log.Println("Starting Finvera Backend...")

	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Connect to Database & Run Migrations
	database.ConnectDB(cfg)

	// 3. Setup Router (Gin)
	r := gin.Default()

	// CORS Config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // Atau gunakan corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(corsConfig))

	// Health Check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "ok",
		})
	})

	// 4. Setup API Routes + Swagger
	router.SetupRouter(r, database.DB, cfg)

	// 5. Start Server
	log.Printf("Server running on port %s", cfg.Port)
	log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
