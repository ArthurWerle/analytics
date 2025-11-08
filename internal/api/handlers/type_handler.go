package handlers

import (
	"analytics/internal/repository"
	"analytics/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TypeHandler struct {
	repo    *repository.TypeRepository
	service *service.TypeService
}

func NewTypeHandler(repo *repository.TypeRepository, service *service.TypeService) *TypeHandler {
	return &TypeHandler{
		repo:    repo,
		service: service,
	}
}

func (h *TypeHandler) GetTypes(c *gin.Context) {
	types, err := h.repo.GetAllTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, types)
}

func (h *TypeHandler) GetAverageByType(c *gin.Context) {
	average, err := h.service.GetAverageByType(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, average)
}
