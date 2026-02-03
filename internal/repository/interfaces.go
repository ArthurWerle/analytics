package repository

import (
	"analytics/internal/domain"
	"context"
)

type TransactionRepositoryInterface interface {
	GetAllTransactions(ctx context.Context) ([]domain.Transaction, error)
}

type CategoryRepositoryInterface interface {
	GetAllCategories(ctx context.Context) ([]domain.Category, error)
}
