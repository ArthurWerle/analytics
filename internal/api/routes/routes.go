package routes

import (
	"analytics/internal/api/handlers"
	"analytics/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, transactionHandler *handlers.TransactionHandler) {
	router.Use(middleware.Logger())

	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	router.GET("/transactions", transactionHandler.GetTransactions)
	router.GET("/average-spend", transactionHandler.GetAverageSpendByCategory)

	router.POST("/v1/query", handlers.GetQueryFromOpenAI)
}
