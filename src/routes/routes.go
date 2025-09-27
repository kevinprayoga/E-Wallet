package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"application-wallet/controllers"
	"application-wallet/middleware"
	"application-wallet/repositories"
	"application-wallet/services"
	"application-wallet/config"
	"net/http"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	// Initialize repositories
	transactionRepo := &repositories.TransactionRepository{DB: db}

	// Initialize services
	transactionService := &services.TransactionService{Repo: transactionRepo}

	// Initialize controllers
	transactionController := &controllers.TransactionController{Service: transactionService}
	authController := &controllers.AuthController{DB: db}

	// Public routes
	router.POST("/login", authController.Login)
	router.GET("/healthz", func(c *gin.Context) {
		if err := config.HealthCheck(db); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
						"status": "unhealthy",
						"error":  err.Error(),
				})
				return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	
	// Protected routes (require authentication)
	protected := router.Group("/transaction")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/topup/:userID", transactionController.TopUp)
		protected.POST("/withdraw/:userID", transactionController.Withdraw)
		protected.POST("/pending", transactionController.UpdatePendingTransaction)
		protected.GET("/all-history/:userID", transactionController.GetAllTransactionHistory)
	}
}
