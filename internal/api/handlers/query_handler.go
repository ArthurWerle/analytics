package handlers

import (
	"analytics/internal/service"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

func GetQueryFromOpenAI(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		c.JSON(400, gin.H{"error": "Failed to read request body"})
		return
	}
	defer c.Request.Body.Close()

	queryService := service.QueryService{}
	response, err := queryService.AnalyzeDatabase(string(body))
	if err != nil {
		log.Printf("Error analyzing database: %v", err)
		c.JSON(400, gin.H{"error": "Error analyzing database"})
		return
	}

	c.JSON(200, response)
}
