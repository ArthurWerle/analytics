package handlers

import (
	"analytics/internal/repository"
	"analytics/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	repo    *repository.TransactionRepository
	service *service.TransactionAnalysisService
}

func NewTransactionHandler(repo *repository.TransactionRepository, service *service.TransactionAnalysisService) *TransactionHandler {
	return &TransactionHandler{
		repo:    repo,
		service: service,
	}
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	transactions, err := h.repo.GetAllTransactions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) GetAverageByCategory(c *gin.Context) {
	average, err := h.service.GetAverageSpendByCategory(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, average)
}
