package domain

import "time"

type Transaction struct {
	ID          int       `db:"id"`
	CategoryID  int       `db:"category_id"`
	Amount      float64   `db:"amount"`
	TypeID      int       `db:"type_id"`
	UpdatedAt   time.Time `db:"updated_at"`
	Date        time.Time `db:"date"`
	CreatedAt   time.Time `db:"created_at"`
	Description string    `db:"description"`
}
