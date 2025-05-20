package repository

import (
	"analytics/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
)

type RecurringTransactionRepository struct {
	db *pgx.Conn
}

func NewRecurringTransactionRepository(db *pgx.Conn) *RecurringTransactionRepository {
	return &RecurringTransactionRepository{db: db}
}

func (r *RecurringTransactionRepository) GetAllRecurringTransactions(ctx context.Context) ([]domain.RecurringTransaction, error) {
	rows, err := r.db.Query(ctx, `
		SELECT 
			id,
			category_id,
			amount,
			type_id,
			updated_at,
			start_date,
			end_date,
			created_at,
			description,
			last_occurrence,
			frequency
		FROM recurring_transactions
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recurringTransactions []domain.RecurringTransaction
	for rows.Next() {
		var recurringTransaction domain.RecurringTransaction
		err := rows.Scan(
			&recurringTransaction.ID,
			&recurringTransaction.CategoryID,
			&recurringTransaction.Amount,
			&recurringTransaction.TypeID,
			&recurringTransaction.UpdatedAt,
			&recurringTransaction.StartDate,
			&recurringTransaction.EndDate,
			&recurringTransaction.CreatedAt,
			&recurringTransaction.Description,
			&recurringTransaction.LastOccurrence,
			&recurringTransaction.Frequency,
		)
		if err != nil {
			return nil, err
		}
		recurringTransactions = append(recurringTransactions, recurringTransaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return recurringTransactions, nil
}
