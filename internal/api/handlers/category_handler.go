package handlers

import (
	"analytics/internal/repository"
	"analytics/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	repo    repository.CategoryRepositoryInterface
	service *service.CategoryService
}

func NewCategoryHandler(repo repository.CategoryRepositoryInterface, service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		repo:    repo,
		service: service,
	}
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.repo.GetAllCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) GetAverageByCategory(c *gin.Context) {
	average, err := h.service.GetAverageByCategory(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, average)
}
