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

type mockTypeRepoForHandler struct {
	types []domain.Type
	err   error
}

func (m *mockTypeRepoForHandler) GetAllTypes(ctx context.Context) ([]domain.Type, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.types, nil
}

type mockTransactionRepoForHandler struct {
	transactions []domain.Transaction
	err          error
}

func (m *mockTransactionRepoForHandler) GetAllTransactions(ctx context.Context) ([]domain.Transaction, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.transactions, nil
}

type mockRecurringRepoForHandler struct {
	recurringTransactions []domain.RecurringTransaction
	err                   error
}

func (m *mockRecurringRepoForHandler) GetAllRecurringTransactions(ctx context.Context) ([]domain.RecurringTransaction, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.recurringTransactions, nil
}

func TestGetAverageByType_Handler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jan2024 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	feb2024 := time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC)

	mockTypeRepo := &mockTypeRepoForHandler{
		types: []domain.Type{
			{ID: 1, Name: domain.Expense},
			{ID: 2, Name: domain.Income},
		},
	}

	mockTransactionRepo := &mockTransactionRepoForHandler{
		transactions: []domain.Transaction{
			{ID: 1, TypeID: 1, Amount: 100.0, Date: jan2024},
			{ID: 2, TypeID: 1, Amount: 200.0, Date: jan2024},
			{ID: 3, TypeID: 2, Amount: 3000.0, Date: feb2024},
		},
	}

	mockRecurringRepo := &mockRecurringRepoForHandler{
		recurringTransactions: []domain.RecurringTransaction{
			{ID: 1, TypeID: 1, Amount: 300.0, StartDate: jan2024, EndDate: sql.NullTime{Valid: false}},
		},
	}

	typeService := service.NewTypeService(mockTypeRepo, mockTransactionRepo, mockRecurringRepo)
	handler := NewTypeHandler(mockTypeRepo, typeService)

	router := gin.New()
	router.GET("/api/v1/types/average", handler.GetAverageByType)

	req, err := http.NewRequest("GET", "/api/v1/types/average", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var results []service.AverageType
	err = json.Unmarshal(w.Body.Bytes(), &results)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results (Expense-Jan, Income-Feb), got %d", len(results))
	}

	for _, result := range results {
		if result.TypeID == 1 && result.Month.Month() == time.January {
			expectedAvg := 200.0
			if result.Average != expectedAvg {
				t.Errorf("Expected expense average %f, got %f", expectedAvg, result.Average)
			}
			if result.TypeName != string(domain.Expense) {
				t.Errorf("Expected type name %s, got %s", domain.Expense, result.TypeName)
			}
		} else if result.TypeID == 2 && result.Month.Month() == time.February {
			expectedAvg := 3000.0
			if result.Average != expectedAvg {
				t.Errorf("Expected income average %f, got %f", expectedAvg, result.Average)
			}
			if result.TypeName != string(domain.Income) {
				t.Errorf("Expected type name %s, got %s", domain.Income, result.TypeName)
			}
		}
	}
}

func TestGetAverageByType_Handler_WithNoTransactions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockTypeRepo := &mockTypeRepoForHandler{
		types: []domain.Type{
			{ID: 1, Name: domain.Expense},
		},
	}

	mockTransactionRepo := &mockTransactionRepoForHandler{
		transactions: []domain.Transaction{},
	}

	mockRecurringRepo := &mockRecurringRepoForHandler{
		recurringTransactions: []domain.RecurringTransaction{},
	}

	typeService := service.NewTypeService(mockTypeRepo, mockTransactionRepo, mockRecurringRepo)
	handler := NewTypeHandler(mockTypeRepo, typeService)

	router := gin.New()
	router.GET("/api/v1/types/average", handler.GetAverageByType)

	req, err := http.NewRequest("GET", "/api/v1/types/average", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var results []service.AverageType
	err = json.Unmarshal(w.Body.Bytes(), &results)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestGetAverageByType_Handler_HandlesNullableFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	may2024 := time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC)

	mockTypeRepo := &mockTypeRepoForHandler{
		types: []domain.Type{
			{
				ID:        1,
				Name:      domain.Expense,
				DeletedAt: sql.NullTime{Valid: false},
			},
		},
	}

	mockTransactionRepo := &mockTransactionRepoForHandler{
		transactions: []domain.Transaction{
			{ID: 1, TypeID: 1, Amount: 500.0, Date: may2024},
		},
	}

	mockRecurringRepo := &mockRecurringRepoForHandler{
		recurringTransactions: []domain.RecurringTransaction{
			{
				ID:             1,
				TypeID:         1,
				Amount:         100.0,
				StartDate:      may2024,
				EndDate:        sql.NullTime{Valid: false},
				LastOccurrence: sql.NullTime{Valid: false},
			},
		},
	}

	typeService := service.NewTypeService(mockTypeRepo, mockTransactionRepo, mockRecurringRepo)
	handler := NewTypeHandler(mockTypeRepo, typeService)

	router := gin.New()
	router.GET("/api/v1/types/average", handler.GetAverageByType)

	req, err := http.NewRequest("GET", "/api/v1/types/average", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var results []service.AverageType
	err = json.Unmarshal(w.Body.Bytes(), &results)
	if err != nil {
		t.Fatalf("Failed to parse response: %v. Body: %s", err, w.Body.String())
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	expectedAvg := 300.0
	if results[0].Average != expectedAvg {
		t.Errorf("Expected average %f, got %f", expectedAvg, results[0].Average)
	}

	if results[0].Month.Month() != time.May || results[0].Month.Year() != 2024 {
		t.Errorf("Expected month May 2024, got %s", results[0].Month)
	}
}
