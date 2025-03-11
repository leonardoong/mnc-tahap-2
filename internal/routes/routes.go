package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leonardoong/e-wallet/internal/handler"
	"github.com/leonardoong/e-wallet/internal/middleware"
	"github.com/leonardoong/e-wallet/internal/service"
)

func SetupRoutes(router *gin.Engine, authService service.IAuthService, transactionService service.ITransactionService) {
	authHandler := handler.AuthHandler{
		AuthService: authService,
	}

	transactionHandler := handler.TransactionHandler{
		TransactionService: transactionService,
	}

	jwtMiddleware := middleware.JWTMiddleware{
		AuthService: authService,
	}

	publicRoutes := router.Group("")
	publicRoutes.POST("/register", authHandler.Register)
	publicRoutes.POST("/login", authHandler.Login)

	protectedRoutes := router.Group("")
	protectedRoutes.Use(jwtMiddleware.AuthRequired())
	protectedRoutes.PUT("/profile", authHandler.UpdateProfile)
	protectedRoutes.POST("/topup", transactionHandler.TopUp)
	protectedRoutes.GET("/topup/:top_up_id", transactionHandler.FindTopUp)
	protectedRoutes.POST("/payment", transactionHandler.Payment)
	protectedRoutes.GET("/payment/:payment_id", transactionHandler.FindPayment)
	protectedRoutes.POST("/transfer", transactionHandler.Transfer)
	protectedRoutes.GET("/transfer/:transfer_id", transactionHandler.FindTransfer)
	protectedRoutes.GET("/transactions", transactionHandler.FindTransactions)
}
