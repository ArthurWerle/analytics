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

	v1 := router.Group("/api/v1")
	v1.POST("/query", handlers.GetQueryFromOpenAI)

	{
		transactions := v1.Group("/transactions")
		transactions.GET("/", transactionHandler.GetTransactions)
	}

	{
		category := v1.Group("/category")
		category.GET("/:category_id/average-spend", transactionHandler.GetAverageSpendByCategory)
	}

	{
		transactionType := v1.Group("/type")
		transactionType.GET("/:type_id/average-spend", transactionHandler.GetAverageSpendByCategory)
	}
}
