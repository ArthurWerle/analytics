package domain

import (
	"time"
)

// Type represents the transaction type (income or expense)
type Type string

const (
	Income  Type = "income"
	Expense Type = "expense"
)

type Transaction struct {
	ID          int     `db:"id"`
	CategoryID  int     `db:"category_id"`
	CreatedById int     `db:"created_by_id"`
	Amount      float64 `db:"amount"`
	Type        Type    `db:"type"`
	Subtype     string       `db:"subtype"`
	UpdatedAt   time.Time    `db:"updated_at"`
	Frequency   string       `db:"frequency"`
	StartDate   *time.Time   `db:"start_date"`
	EndDate     *time.Time   `db:"end_date"`
	Date        time.Time    `db:"date"`
	CreatedAt   time.Time    `db:"created_at"`
	Description string       `db:"description"`
	IsRecurring bool         `db:"is_recurring"`
}
