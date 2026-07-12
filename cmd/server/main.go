package main

import (
	"context"
	"finvera-be/internal/config"
	"finvera-be/internal/database"
	"finvera-be/internal/middleware"
	"finvera-be/internal/router"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "finvera-be/docs" // swagger docs

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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting Finvera Backend...")

	// 1. Load Config (validates secrets, origins, etc.)
	cfg := config.LoadConfig()

	// 2. Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// 3. Connect to Database (returns *gorm.DB — no more global-only access)
	db := database.ConnectDB(cfg)

	// 4. Setup Gin engine with custom middleware only (no default logger+recovery in production)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 5. Security Headers middleware (applies to all routes)
	r.Use(middleware.SecurityHeaders())

	// 6. CORS — origins from config, not AllowAllOrigins
	corsConfig := cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	// 7. Health Check
	r.GET("/ping", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"message": "database unavailable",
				"status":  "error",
				"env":     cfg.AppEnv,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "ok",
			"env":     cfg.AppEnv,
		})
	})

	// 8. Setup all API routes (DI happens inside SetupRouter and returns services for cron)
	cronService := router.SetupRouter(r, db, cfg)

	// 9. Start Cron Service
	cronService.Start()
	defer cronService.Stop()

	// 10. Start Server with Graceful Shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		slog.Info("Server running", "port", cfg.Port, "env", cfg.AppEnv)
		if !cfg.IsProduction() {
			slog.Info("Swagger UI: http://localhost:" + cfg.Port + "/swagger/index.html")
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	// 5 seconds timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting gracefully")
}
