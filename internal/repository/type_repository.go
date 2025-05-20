package repository

import (
	"analytics/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
)

type TypeRepository struct {
	db *pgx.Conn
}

func NewTypeRepository(db *pgx.Conn) *TypeRepository {
	return &TypeRepository{db: db}
}

func (r *TypeRepository) GetAllTypes(ctx context.Context) ([]domain.Type, error) {
	rows, err := r.db.Query(ctx, `
		SELECT 
			id,
			updated_at,
			created_at,
			deleted_at,
			name,
			description
		FROM types
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []domain.Type
	for rows.Next() {
		var t domain.Type
		err := rows.Scan(
			&t.ID,
			&t.UpdatedAt,
			&t.CreatedAt,
			&t.DeletedAt,
			&t.Name,
			&t.Description,
		)
		if err != nil {
			return nil, err
		}
		types = append(types, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return types, nil
}
