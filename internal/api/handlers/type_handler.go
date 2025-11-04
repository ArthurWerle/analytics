package handlers

import (
	"analytics/internal/repository"
	"analytics/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionType string

const (
	Expense TransactionType = "expense"
	Income  TransactionType = "income"
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

func (h *TypeHandler) GetAverageSpendByType(c *gin.Context) {
	averageSpend, err := h.service.GetAverageSpendByType(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, averageSpend)
}
