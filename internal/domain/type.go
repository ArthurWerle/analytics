package domain

import "time"

type Type struct {
	ID          int       `db:"id"`
	UpdatedAt   time.Time `db:"updated_at"`
	CreatedAt   time.Time `db:"created_at"`
	DeletedAt   time.Time `db:"deleted_at"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
}
