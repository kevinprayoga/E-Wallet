package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"application-wallet/controllers"
	"application-wallet/middleware"
	"application-wallet/repositories"
	"application-wallet/services"
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
