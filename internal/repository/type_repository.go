package repository

import (
	"analytics/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TypeRepository struct {
	db *pgxpool.Pool
}

func NewTypeRepository(db *pgxpool.Pool) *TypeRepository {
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
