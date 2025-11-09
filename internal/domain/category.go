package domain

import (
	"database/sql"
	"time"
)

type Category struct {
	ID          int          `db:"id"`
	UpdatedAt   time.Time    `db:"updated_at"`
	CreatedAt   time.Time    `db:"created_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
	Name        string       `db:"name"`
	Description string       `db:"description"`
	Color       string       `db:"color"`
}
