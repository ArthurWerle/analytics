package service

import (
	"analytics/internal/domain"
	"context"
	"database/sql"
	"testing"
	"time"
)

type mockTypeRepository struct {
	types []domain.Type
	err   error
}

func (m *mockTypeRepository) GetAllTypes(ctx context.Context) ([]domain.Type, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.types, nil
}

type mockTransactionRepository struct {
	transactions []domain.Transaction
	err          error
}

func (m *mockTransactionRepository) GetAllTransactions(ctx context.Context) ([]domain.Transaction, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.transactions, nil
}

type mockRecurringTransactionRepository struct {
	recurringTransactions []domain.RecurringTransaction
	err                   error
}

func (m *mockRecurringTransactionRepository) GetAllRecurringTransactions(ctx context.Context) ([]domain.RecurringTransaction, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.recurringTransactions, nil
}

func TestGetAverageByType_WithOnlyRegularTransactions(t *testing.T) {
	ctx := context.Background()
	jan2024 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	feb2024 := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)

	mockTypeRepo := &mockTypeRepository{
		types: []domain.Type{
			{ID: 1, Name: domain.Expense},
			{ID: 2, Name: domain.Income},
		},
	}

	mockTransactionRepo := &mockTransactionRepository{
		transactions: []domain.Transaction{
			{ID: 1, TypeID: 1, Amount: 100.0, Date: jan2024},
			{ID: 2, TypeID: 1, Amount: 200.0, Date: jan2024},
			{ID: 3, TypeID: 2, Amount: 1000.0, Date: feb2024},
			{ID: 4, TypeID: 2, Amount: 2000.0, Date: feb2024},
		},
	}

	mockRecurringRepo := &mockRecurringTransactionRepository{
		recurringTransactions: []domain.RecurringTransaction{},
	}

	service := NewTypeService(mockTypeRepo, mockTransactionRepo, mockRecurringRepo)

	results, err := service.GetAverageByType(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results (one per type per month), got %d", len(results))
	}

	for _, result := range results {
		if result.TypeID == 1 && result.Month.Month() == time.January {
			expectedAvg := 150.0
			if result.Average != expectedAvg {
				t.Errorf("Expected expense average in January %f, got %f", expectedAvg, result.Average)
			}
			if result.TypeName != string(domain.Expense) {
				t.Errorf("Expected type name %s, got %s", domain.Expense, result.TypeName)
			}
		} else if result.TypeID == 2 && result.Month.Month() == time.February {
			expectedAvg := 1500.0
			if result.Average != expectedAvg {
				t.Errorf("Expected income average in February %f, got %f", expectedAvg, result.Average)
			}
			if result.TypeName != string(domain.Income) {
				t.Errorf("Expected type name %s, got %s", domain.Income, result.TypeName)
			}
		}
	}
}

func TestGetAverageByType_WithBothRegularAndRecurringTransactions(t *testing.T) {
	ctx := context.Background()
	jan2024 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	mockTypeRepo := &mockTypeRepository{
		types: []domain.Type{
			{ID: 1, Name: domain.Expense},
		},
	}

	mockTransactionRepo := &mockTransactionRepository{
		transactions: []domain.Transaction{
			{ID: 1, TypeID: 1, Amount: 100.0, Date: jan2024},
			{ID: 2, TypeID: 1, Amount: 200.0, Date: jan2024},
		},
	}

	mockRecurringRepo := &mockRecurringTransactionRepository{
		recurringTransactions: []domain.RecurringTransaction{
			{
				ID:        1,
				TypeID:    1,
				Amount:    300.0,
				StartDate: jan2024,
				EndDate:   sql.NullTime{Valid: false},
			},
		},
	}

	service := NewTypeService(mockTypeRepo, mockTransactionRepo, mockRecurringRepo)

	results, err := service.GetAverageByType(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	expectedAvg := 200.0
	if results[0].Average != expectedAvg {
		t.Errorf("Expected average %f, got %f", expectedAvg, results[0].Average)
	}

	if results[0].TypeName != string(domain.Expense) {
		t.Errorf("Expected type name %s, got %s", domain.Expense, results[0].TypeName)
	}

	if results[0].Month.Month() != time.January || results[0].Month.Year() != 2024 {
		t.Errorf("Expected month January 2024, got %s", results[0].Month)
	}
}

func TestGetAverageByType_WithNoTransactions(t *testing.T) {
	ctx := context.Background()

	mockTypeRepo := &mockTypeRepository{
		types: []domain.Type{
			{ID: 1, Name: domain.Expense},
		},
	}

	mockTransactionRepo := &mockTransactionRepository{
		transactions: []domain.Transaction{},
	}

	mockRecurringRepo := &mockRecurringTransactionRepository{
		recurringTransactions: []domain.RecurringTransaction{},
	}

	service := NewTypeService(mockTypeRepo, mockTransactionRepo, mockRecurringRepo)

	results, err := service.GetAverageByType(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestGetAverageByType_WithSingleTransaction(t *testing.T) {
	ctx := context.Background()
	march2024 := time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)

	mockTypeRepo := &mockTypeRepository{
		types: []domain.Type{
			{ID: 1, Name: domain.Expense},
		},
	}

	mockTransactionRepo := &mockTransactionRepository{
		transactions: []domain.Transaction{
			{ID: 1, TypeID: 1, Amount: 250.50, Date: march2024},
		},
	}

	mockRecurringRepo := &mockRecurringTransactionRepository{
		recurringTransactions: []domain.RecurringTransaction{},
	}

	service := NewTypeService(mockTypeRepo, mockTransactionRepo, mockRecurringRepo)

	results, err := service.GetAverageByType(ctx)

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

	if results[0].Month.Month() != time.March || results[0].Month.Year() != 2024 {
		t.Errorf("Expected month March 2024, got %s", results[0].Month)
	}
}

func TestGetAverageByType_VerifiesTypeNamesAreMapped(t *testing.T) {
	ctx := context.Background()
	april2024 := time.Date(2024, 4, 5, 0, 0, 0, 0, time.UTC)

	mockTypeRepo := &mockTypeRepository{
		types: []domain.Type{
			{ID: 1, Name: domain.Expense},
			{ID: 2, Name: domain.Income},
		},
	}

	mockTransactionRepo := &mockTransactionRepository{
		transactions: []domain.Transaction{
			{ID: 1, TypeID: 1, Amount: 100.0, Date: april2024},
			{ID: 2, TypeID: 2, Amount: 500.0, Date: april2024},
		},
	}

	mockRecurringRepo := &mockRecurringTransactionRepository{
		recurringTransactions: []domain.RecurringTransaction{},
	}

	service := NewTypeService(mockTypeRepo, mockTransactionRepo, mockRecurringRepo)

	results, err := service.GetAverageByType(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	typeNames := make(map[int]string)
	for _, result := range results {
		typeNames[result.TypeID] = result.TypeName
	}

	if typeNames[1] != string(domain.Expense) {
		t.Errorf("Expected type name for ID 1 to be %s, got %s", domain.Expense, typeNames[1])
	}

	if typeNames[2] != string(domain.Income) {
		t.Errorf("Expected type name for ID 2 to be %s, got %s", domain.Income, typeNames[2])
	}
}
