package handlers

import (
	"analytics/internal/service"
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
	average, err := h.service.GetAverageByType(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, average)
}
