package repository

import (
	"analytics/internal/domain"
	"context"
)

// TypeRepositoryInterface defines the contract for type repository operations
type TypeRepositoryInterface interface {
	GetAllTypes(ctx context.Context) ([]domain.Type, error)
}

// TransactionRepositoryInterface defines the contract for transaction repository operations
type TransactionRepositoryInterface interface {
	GetAllTransactions(ctx context.Context) ([]domain.Transaction, error)
}

// RecurringTransactionRepositoryInterface defines the contract for recurring transaction repository operations
type RecurringTransactionRepositoryInterface interface {
	GetAllRecurringTransactions(ctx context.Context) ([]domain.RecurringTransaction, error)
}
