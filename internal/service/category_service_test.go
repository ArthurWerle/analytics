package service

import (
	"analytics/internal/domain"
	"context"
	"database/sql"
	"testing"
	"time"
)

type mockCategoryRepository struct {
	categories []domain.Category
	err        error
}

func (m *mockCategoryRepository) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.categories, nil
}

func TestGetAverageByCategory_WithOnlyRegularTransactions(t *testing.T) {
	ctx := context.Background()
	jan2024 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	feb2024 := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)

	mockCategoryRepo := &mockCategoryRepository{
		categories: []domain.Category{
			{ID: 1, Name: "Groceries"},
			{ID: 2, Name: "Transportation"},
		},
	}

	mockTransactionRepo := &mockTransactionRepository{
		transactions: []domain.Transaction{
			{ID: 1, CategoryID: 1, Amount: 100.0, Date: jan2024},
			{ID: 2, CategoryID: 1, Amount: 200.0, Date: jan2024},
			{ID: 3, CategoryID: 1, Amount: 300.0, Date: feb2024},
			{ID: 4, CategoryID: 2, Amount: 1000.0, Date: jan2024},
			{ID: 5, CategoryID: 2, Amount: 2000.0, Date: feb2024},
		},
	}

	mockRecurringRepo := &mockRecurringTransactionRepository{
		recurringTransactions: []domain.RecurringTransaction{},
	}

	service := NewCategoryService(mockCategoryRepo, mockTransactionRepo, mockRecurringRepo)

	results, err := service.GetAverageByCategory(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results (one per category), got %d", len(results))
	}

	for _, result := range results {
		if result.CategoryID == 1 {
			expectedAvg := (300.0 + 300.0) / 2
			if result.Average != expectedAvg {
				t.Errorf("Expected Groceries average %f, got %f", expectedAvg, result.Average)
			}
			if result.CategoryName != "Groceries" {
				t.Errorf("Expected category name Groceries, got %s", result.CategoryName)
			}
		} else if result.CategoryID == 2 {
			expectedAvg := (1000.0 + 2000.0) / 2
			if result.Average != expectedAvg {
				t.Errorf("Expected Transportation average %f, got %f", expectedAvg, result.Average)
			}
			if result.CategoryName != "Transportation" {
				t.Errorf("Expected category name Transportation, got %s", result.CategoryName)
			}
		}
	}
}

func TestGetAverageByCategory_WithBothRegularAndRecurringTransactions(t *testing.T) {
	ctx := context.Background()
	jan2024 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	mockCategoryRepo := &mockCategoryRepository{
		categories: []domain.Category{
			{ID: 1, Name: "Groceries"},
		},
	}

	mockTransactionRepo := &mockTransactionRepository{
		transactions: []domain.Transaction{
			{ID: 1, CategoryID: 1, Amount: 100.0, Date: jan2024},
			{ID: 2, CategoryID: 1, Amount: 200.0, Date: jan2024},
		},
	}

	mockRecurringRepo := &mockRecurringTransactionRepository{
		recurringTransactions: []domain.RecurringTransaction{
			{
				ID:         1,
				CategoryID: 1,
				Amount:     300.0,
				StartDate:  jan2024,
				EndDate:    sql.NullTime{Valid: false},
			},
		},
	}

	service := NewCategoryService(mockCategoryRepo, mockTransactionRepo, mockRecurringRepo)

	results, err := service.GetAverageByCategory(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	expectedAvg := 600.0
	if results[0].Average != expectedAvg {
		t.Errorf("Expected average %f, got %f", expectedAvg, results[0].Average)
	}

	if results[0].CategoryName != "Groceries" {
		t.Errorf("Expected category name Groceries, got %s", results[0].CategoryName)
	}
}

func TestGetAverageByCategory_WithNoTransactions(t *testing.T) {
	ctx := context.Background()

	mockCategoryRepo := &mockCategoryRepository{
		categories: []domain.Category{
			{ID: 1, Name: "Groceries"},
		},
	}

	mockTransactionRepo := &mockTransactionRepository{
		transactions: []domain.Transaction{},
	}

	mockRecurringRepo := &mockRecurringTransactionRepository{
		recurringTransactions: []domain.RecurringTransaction{},
	}

	service := NewCategoryService(mockCategoryRepo, mockTransactionRepo, mockRecurringRepo)

	results, err := service.GetAverageByCategory(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestGetAverageByCategory_WithSingleTransaction(t *testing.T) {
	ctx := context.Background()
	march2024 := time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)

	mockCategoryRepo := &mockCategoryRepository{
		categories: []domain.Category{
			{ID: 1, Name: "Groceries"},
		},
	}

	mockTransactionRepo := &mockTransactionRepository{
		transactions: []domain.Transaction{
			{ID: 1, CategoryID: 1, Amount: 250.50, Date: march2024},
		},
	}

	mockRecurringRepo := &mockRecurringTransactionRepository{
		recurringTransactions: []domain.RecurringTransaction{},
	}

	service := NewCategoryService(mockCategoryRepo, mockTransactionRepo, mockRecurringRepo)

	results, err := service.GetAverageByCategory(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	expectedAvg := 250.50
	if results[0].Average != expectedAvg {
		t.Errorf("Expected average %f, got %f", expectedAvg, results[0].Average)
	}
}

func TestGetAverageByCategory_VerifiesCategoryNamesAreMapped(t *testing.T) {
	ctx := context.Background()
	april2024 := time.Date(2024, 4, 5, 0, 0, 0, 0, time.UTC)

	mockCategoryRepo := &mockCategoryRepository{
		categories: []domain.Category{
			{ID: 1, Name: "Groceries"},
			{ID: 2, Name: "Transportation"},
		},
	}

	mockTransactionRepo := &mockTransactionRepository{
		transactions: []domain.Transaction{
			{ID: 1, CategoryID: 1, Amount: 100.0, Date: april2024},
			{ID: 2, CategoryID: 2, Amount: 500.0, Date: april2024},
		},
	}

	mockRecurringRepo := &mockRecurringTransactionRepository{
		recurringTransactions: []domain.RecurringTransaction{},
	}

	service := NewCategoryService(mockCategoryRepo, mockTransactionRepo, mockRecurringRepo)

	results, err := service.GetAverageByCategory(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	categoryNames := make(map[int]string)
	for _, result := range results {
		categoryNames[result.CategoryID] = result.CategoryName
	}

	if categoryNames[1] != "Groceries" {
		t.Errorf("Expected category name for ID 1 to be Groceries, got %s", categoryNames[1])
	}

	if categoryNames[2] != "Transportation" {
		t.Errorf("Expected category name for ID 2 to be Transportation, got %s", categoryNames[2])
	}
}
