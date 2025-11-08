package domain

import (
	"database/sql"
	"time"
)

type TransactionType string

const (
	Expense TransactionType = "expense"
	Income  TransactionType = "income"
)

type Type struct {
	ID          int             `db:"id"`
	UpdatedAt   time.Time       `db:"updated_at"`
	CreatedAt   time.Time       `db:"created_at"`
	DeletedAt   sql.NullTime    `db:"deleted_at"`
	Name        TransactionType `db:"name"`
	Description string          `db:"description"`
}
