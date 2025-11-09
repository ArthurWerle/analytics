package repository

import (
	"analytics/internal/domain"
	"context"
)

type TypeRepositoryInterface interface {
	GetAllTypes(ctx context.Context) ([]domain.Type, error)
}

type TransactionRepositoryInterface interface {
	GetAllTransactions(ctx context.Context) ([]domain.Transaction, error)
}

type RecurringTransactionRepositoryInterface interface {
	GetAllRecurringTransactions(ctx context.Context) ([]domain.RecurringTransaction, error)
}

type CategoryRepositoryInterface interface {
	GetAllCategories(ctx context.Context) ([]domain.Category, error)
}
