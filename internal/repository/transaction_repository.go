package repository

import (
	"analytics/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) GetAllTransactions(ctx context.Context) ([]domain.Transaction, error) {
	rows, err := r.db.Query(ctx, `
		SELECT 
			id,
			category_id,
			amount,
			type,
			updated_at,
			date,
			created_at,
		    start_date,
		    end_date,
			description
		FROM transactions
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.CategoryID,
			&transaction.Amount,
			&transaction.Type,
			&transaction.UpdatedAt,
			&transaction.Date,
			&transaction.CreatedAt,
			&transaction.StartDate,
			&transaction.EndDate,
			&transaction.Description,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
