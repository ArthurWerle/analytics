package routes

import (
	"analytics/internal/api/handlers"
	"analytics/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, transactionHandler *handlers.TransactionHandler, typeHandler *handlers.TypeHandler) {
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
		categories := v1.Group("/categories")
		categories.GET("/:category_id/average", transactionHandler.GetAverageByCategory)
	}

	{
		types := v1.Group("/types")
		types.GET("/average", typeHandler.GetAverageByType)
		//types.GET("/:type_id/average", typeHandler.GetAverageByType) // For future use
	}
}
