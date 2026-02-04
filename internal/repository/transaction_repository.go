package repository

import (
	"analytics/internal/domain"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) GetAllTransactions(ctx context.Context) ([]domain.Transaction, error) {
	log.Printf("[TransactionRepository.GetAllTransactions] Executing query to fetch all transactions")

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
		log.Printf("[TransactionRepository.GetAllTransactions] ERROR: Query failed: %v", err)
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var transactions []domain.Transaction
	rowCount := 0
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
			log.Printf("[TransactionRepository.GetAllTransactions] ERROR: Failed to scan row %d: %v", rowCount, err)
			return nil, fmt.Errorf("failed to scan row %d: %w", rowCount, err)
		}
		transactions = append(transactions, transaction)
		rowCount++
	}

	if err := rows.Err(); err != nil {
		log.Printf("[TransactionRepository.GetAllTransactions] ERROR: Row iteration error: %v", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	log.Printf("[TransactionRepository.GetAllTransactions] Successfully fetched %d transactions", len(transactions))

	return transactions, nil
}
