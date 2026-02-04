package handlers

import (
	"analytics/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TypeHandler struct {
	service *service.TypeService
}

func NewTypeHandler(service *service.TypeService) *TypeHandler {
	return &TypeHandler{
		service: service,
	}
}

func (h *TypeHandler) GetAverageByType(c *gin.Context) {
	log.Printf("[TypeHandler.GetAverageByType] Starting request from %s", c.ClientIP())

	average, err := h.service.GetAverageByType(c.Request.Context())
	if err != nil {
		log.Printf("[TypeHandler.GetAverageByType] ERROR: Failed to get average by type: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[TypeHandler.GetAverageByType] Successfully retrieved %d type averages", len(average))
	c.JSON(http.StatusOK, average)
}
