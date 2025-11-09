package handlers

import (
	"analytics/internal/domain"
	"analytics/internal/service"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type mockCategoryRepoForHandler struct {
	categories []domain.Category
	err        error
}

func (m *mockCategoryRepoForHandler) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.categories, nil
}

func TestGetAverageByCategory_Handler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jan2024 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	feb2024 := time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC)

	mockCategoryRepo := &mockCategoryRepoForHandler{
		categories: []domain.Category{
			{ID: 1, Name: "Groceries"},
			{ID: 2, Name: "Transportation"},
		},
	}

	mockTransactionRepo := &mockTransactionRepoForHandler{
		transactions: []domain.Transaction{
			{ID: 1, CategoryID: 1, Amount: 100.0, Date: jan2024},
			{ID: 2, CategoryID: 1, Amount: 200.0, Date: jan2024},
			{ID: 3, CategoryID: 2, Amount: 3000.0, Date: feb2024},
		},
	}

	mockRecurringRepo := &mockRecurringRepoForHandler{
		recurringTransactions: []domain.RecurringTransaction{
			{ID: 1, CategoryID: 1, Amount: 300.0, StartDate: jan2024, EndDate: sql.NullTime{Valid: false}},
		},
	}

	categoryService := service.NewCategoryService(mockCategoryRepo, mockTransactionRepo, mockRecurringRepo)
	handler := NewCategoryHandler(mockCategoryRepo, categoryService)

	router := gin.New()
	router.GET("/api/v1/categories/average", handler.GetAverageByCategory)

	req, err := http.NewRequest("GET", "/api/v1/categories/average", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var results []service.AverageCategory
	err = json.Unmarshal(w.Body.Bytes(), &results)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results (one per category), got %d", len(results))
	}

	for _, result := range results {
		if result.CategoryID == 1 {
			expectedAvg := 600.0
			if result.Average != expectedAvg {
				t.Errorf("Expected Groceries average %f, got %f", expectedAvg, result.Average)
			}
			if result.CategoryName != "Groceries" {
				t.Errorf("Expected category name Groceries, got %s", result.CategoryName)
			}
		} else if result.CategoryID == 2 {
			expectedAvg := 3000.0
			if result.Average != expectedAvg {
				t.Errorf("Expected Transportation average %f, got %f", expectedAvg, result.Average)
			}
			if result.CategoryName != "Transportation" {
				t.Errorf("Expected category name Transportation, got %s", result.CategoryName)
			}
		}
	}
}

func TestGetAverageByCategory_Handler_WithNoTransactions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCategoryRepo := &mockCategoryRepoForHandler{
		categories: []domain.Category{
			{ID: 1, Name: "Groceries"},
		},
	}

	mockTransactionRepo := &mockTransactionRepoForHandler{
		transactions: []domain.Transaction{},
	}

	mockRecurringRepo := &mockRecurringRepoForHandler{
		recurringTransactions: []domain.RecurringTransaction{},
	}

	categoryService := service.NewCategoryService(mockCategoryRepo, mockTransactionRepo, mockRecurringRepo)
	handler := NewCategoryHandler(mockCategoryRepo, categoryService)

	router := gin.New()
	router.GET("/api/v1/categories/average", handler.GetAverageByCategory)

	req, err := http.NewRequest("GET", "/api/v1/categories/average", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var results []service.AverageCategory
	err = json.Unmarshal(w.Body.Bytes(), &results)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestGetAverageByCategory_Handler_HandlesNullableFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	may2024 := time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC)

	mockCategoryRepo := &mockCategoryRepoForHandler{
		categories: []domain.Category{
			{
				ID:        1,
				Name:      "Groceries",
				DeletedAt: sql.NullTime{Valid: false},
			},
		},
	}

	mockTransactionRepo := &mockTransactionRepoForHandler{
		transactions: []domain.Transaction{
			{ID: 1, CategoryID: 1, Amount: 500.0, Date: may2024},
		},
	}

	mockRecurringRepo := &mockRecurringRepoForHandler{
		recurringTransactions: []domain.RecurringTransaction{
			{
				ID:             1,
				CategoryID:     1,
				Amount:         100.0,
				StartDate:      may2024,
				EndDate:        sql.NullTime{Valid: false},
				LastOccurrence: sql.NullTime{Valid: false},
			},
		},
	}

	categoryService := service.NewCategoryService(mockCategoryRepo, mockTransactionRepo, mockRecurringRepo)
	handler := NewCategoryHandler(mockCategoryRepo, categoryService)

	router := gin.New()
	router.GET("/api/v1/categories/average", handler.GetAverageByCategory)

	req, err := http.NewRequest("GET", "/api/v1/categories/average", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var results []service.AverageCategory
	err = json.Unmarshal(w.Body.Bytes(), &results)
	if err != nil {
		t.Fatalf("Failed to parse response: %v. Body: %s", err, w.Body.String())
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	expectedAvg := 600.0
	if results[0].Average != expectedAvg {
		t.Errorf("Expected average %f, got %f", expectedAvg, results[0].Average)
	}
}
