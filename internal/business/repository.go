package business

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, b *Business) error {
	query := `
		INSERT INTO businesses (
			firebase_uid,
			business_name,
			owner_name,
			phone,
			email,
			address,
			logo_url
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		b.FirebaseUID,
		b.BusinessName,
		b.OwnerName,
		b.Phone,
		b.Email,
		b.Address,
		b.LogoURL,
	).Scan(
		&b.ID,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create business: %w", err)
	}

	return nil
}
