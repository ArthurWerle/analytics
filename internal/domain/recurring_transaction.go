package domain

import "time"

type RecurringTransaction struct {
	ID             int       `db:"id"`
	CategoryID     int       `db:"category_id"`
	Amount         float64   `db:"amount"`
	TypeID         int       `db:"type_id"`
	UpdatedAt      time.Time `db:"updated_at"`
	StartDate      time.Time `db:"start_date"`
	EndDate        time.Time `db:"end_date"`
	CreatedAt      time.Time `db:"created_at"`
	Description    string    `db:"description"`
	LastOccurrence time.Time `db:"last_occurrence"`
	Frequency      string    `db:"frequency"`
}
