package service

import (
	"analytics/internal/domain"
	"context"
	"database/sql"
	"testing"
	"time"
)

// Mock repositories for testing
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
	// Arrange
	ctx := context.Background()

	mockTypeRepo := &mockTypeRepository{
		types: []domain.Type{
			{ID: 1, Name: domain.Expense},
			{ID: 2, Name: domain.Income},
		},
	}

	mockTransactionRepo := &mockTransactionRepository{
		transactions: []domain.Transaction{
			{ID: 1, TypeID: 1, Amount: 100.0, Date: time.Now()},
			{ID: 2, TypeID: 1, Amount: 200.0, Date: time.Now()},
			{ID: 3, TypeID: 2, Amount: 1000.0, Date: time.Now()},
			{ID: 4, TypeID: 2, Amount: 2000.0, Date: time.Now()},
		},
	}

	mockRecurringRepo := &mockRecurringTransactionRepository{
		recurringTransactions: []domain.RecurringTransaction{},
	}

	service := NewTypeService(mockTypeRepo, mockTransactionRepo, mockRecurringRepo)

	// Act
	results, err := service.GetAverageByType(ctx)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Find expense and income averages
	var expenseAvg, incomeAvg float64
	for _, result := range results {
		if result.TypeID == 1 {
			expenseAvg = result.Average
		} else if result.TypeID == 2 {
			incomeAvg = result.Average
		}
	}

	// Expected: Expense average = (100 + 200) / 2 = 150
	expectedExpenseAvg := 150.0
	if expenseAvg != expectedExpenseAvg {
		t.Errorf("Expected expense average %f, got %f", expectedExpenseAvg, expenseAvg)
	}

	// Expected: Income average = (1000 + 2000) / 2 = 1500
	expectedIncomeAvg := 1500.0
	if incomeAvg != expectedIncomeAvg {
		t.Errorf("Expected income average %f, got %f", expectedIncomeAvg, incomeAvg)
	}
}

func TestGetAverageByType_WithBothRegularAndRecurringTransactions(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockTypeRepo := &mockTypeRepository{
		types: []domain.Type{
			{ID: 1, Name: domain.Expense},
		},
	}

	mockTransactionRepo := &mockTransactionRepository{
		transactions: []domain.Transaction{
			{ID: 1, TypeID: 1, Amount: 100.0, Date: time.Now()},
			{ID: 2, TypeID: 1, Amount: 200.0, Date: time.Now()},
		},
	}

	mockRecurringRepo := &mockRecurringTransactionRepository{
		recurringTransactions: []domain.RecurringTransaction{
			{
				ID:        1,
				TypeID:    1,
				Amount:    300.0,
				StartDate: time.Now(),
				EndDate:   sql.NullTime{Valid: false}, // NULL end date
			},
		},
	}

	service := NewTypeService(mockTypeRepo, mockTransactionRepo, mockRecurringRepo)

	// Act
	results, err := service.GetAverageByType(ctx)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	// Expected: Average = (100 + 200 + 300) / 3 = 200
	expectedAvg := 200.0
	if results[0].Average != expectedAvg {
		t.Errorf("Expected average %f, got %f", expectedAvg, results[0].Average)
	}

	// Verify type name is correct
	if results[0].TypeName != string(domain.Expense) {
		t.Errorf("Expected type name %s, got %s", domain.Expense, results[0].TypeName)
	}
}

func TestGetAverageByType_WithNoTransactions(t *testing.T) {
	// Arrange
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

	// Act
	results, err := service.GetAverageByType(ctx)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return empty results when there are no transactions
	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestGetAverageByType_WithSingleTransaction(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockTypeRepo := &mockTypeRepository{
		types: []domain.Type{
			{ID: 1, Name: domain.Expense},
		},
	}

	mockTransactionRepo := &mockTransactionRepository{
		transactions: []domain.Transaction{
			{ID: 1, TypeID: 1, Amount: 250.50, Date: time.Now()},
		},
	}

	mockRecurringRepo := &mockRecurringTransactionRepository{
		recurringTransactions: []domain.RecurringTransaction{},
	}

	service := NewTypeService(mockTypeRepo, mockTransactionRepo, mockRecurringRepo)

	// Act
	results, err := service.GetAverageByType(ctx)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	// Average of single transaction should be the transaction amount
	expectedAvg := 250.50
	if results[0].Average != expectedAvg {
		t.Errorf("Expected average %f, got %f", expectedAvg, results[0].Average)
	}
}

func TestGetAverageByType_VerifiesTypeNamesAreMapped(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockTypeRepo := &mockTypeRepository{
		types: []domain.Type{
			{ID: 1, Name: domain.Expense},
			{ID: 2, Name: domain.Income},
		},
	}

	mockTransactionRepo := &mockTransactionRepository{
		transactions: []domain.Transaction{
			{ID: 1, TypeID: 1, Amount: 100.0, Date: time.Now()},
			{ID: 2, TypeID: 2, Amount: 500.0, Date: time.Now()},
		},
	}

	mockRecurringRepo := &mockRecurringTransactionRepository{
		recurringTransactions: []domain.RecurringTransaction{},
	}

	service := NewTypeService(mockTypeRepo, mockTransactionRepo, mockRecurringRepo)

	// Act
	results, err := service.GetAverageByType(ctx)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify that type names are correctly mapped
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
