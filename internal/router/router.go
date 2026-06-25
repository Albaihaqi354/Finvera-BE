package router

import (
	"finvera-be/internal/config"
	"finvera-be/internal/handler"
	"finvera-be/internal/middleware"
	"finvera-be/internal/repository"
	"finvera-be/internal/service"

	_ "finvera-be/docs" // swagger docs

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// ─── Swagger UI ──────────────────────────────────────────────────────
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ─── Setup Dependencies ───────────────────────────────────────────────
	// Repositories
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	scheduledRepo := repository.NewScheduledTransactionRepository(db)

	// Services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(accountRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo, categoryRepo, tagRepo)
	scheduledService := service.NewScheduledTransactionService(scheduledRepo, accountRepo, categoryRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService, cfg)
	userHandler := handler.NewUserHandler(userService)
	accountHandler := handler.NewAccountHandler(accountService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	tagHandler := handler.NewTagHandler(tagService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	scheduledHandler := handler.NewScheduledTransactionHandler(scheduledService)

	// ─── API v1 Routes ────────────────────────────────────────────────────
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
		}

		v1.GET("/preset-categories", categoryHandler.GetPresetCategories)

		// Protected Routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Users
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetProfile)
				users.PUT("/me", userHandler.UpdateProfile)
			}

			// Accounts
			accounts := protected.Group("/accounts")
			{
				accounts.POST("", accountHandler.CreateAccount)
				accounts.GET("", accountHandler.GetAccounts)
				accounts.GET("/:id", accountHandler.GetAccountByID)
				accounts.PUT("/:id", accountHandler.UpdateAccount)
				accounts.DELETE("/:id", accountHandler.DeleteAccount)
			}

			// Categories
			categories := protected.Group("/categories")
			{
				categories.POST("", categoryHandler.CreateCategory)
				categories.GET("", categoryHandler.GetCategories)
				categories.GET("/:id", categoryHandler.GetCategoryByID)
				categories.PUT("/:id", categoryHandler.UpdateCategory)
				categories.DELETE("/:id", categoryHandler.DeleteCategory)
			}

			// Tags
			tags := protected.Group("/tags")
			{
				tags.POST("", tagHandler.CreateTag)
				tags.GET("", tagHandler.GetTags)
				tags.GET("/:id", tagHandler.GetTagByID)
				tags.PUT("/:id", tagHandler.UpdateTag)
				tags.DELETE("/:id", tagHandler.DeleteTag)
			}

			// Transactions
			transactions := protected.Group("/transactions")
			{
				transactions.POST("", transactionHandler.CreateTransaction)
				transactions.GET("", transactionHandler.GetTransactions)
				transactions.GET("/:id", transactionHandler.GetTransactionByID)
				transactions.PUT("/:id", transactionHandler.UpdateTransaction)
				transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
			}

			// Scheduled Transactions
			scheduled := protected.Group("/scheduled")
			{
				scheduled.POST("", scheduledHandler.CreateScheduled)
				scheduled.GET("", scheduledHandler.GetScheduleds)
				scheduled.GET("/:id", scheduledHandler.GetScheduledByID)
				scheduled.PUT("/:id", scheduledHandler.UpdateScheduled)
				scheduled.DELETE("/:id", scheduledHandler.DeleteScheduled)
			}
		}
	}
}
