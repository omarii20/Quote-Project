package business

import (
	"context"
	"database/sql"
	"errors"
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

// Create inserts a new business into the database and returns the created business with its ID and timestamps.
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
	).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create business: %w", err)
	}

	return nil
}

// GetByID retrieves a business by its ID from the database.
func (r *Repository) GetByID(ctx context.Context, id int64) (*Business, error) {

	query := `
		SELECT
			id,
			firebase_uid,
			business_name,
			owner_name,
			phone,
			email,
			address,
			logo_url,
			created_at,
			updated_at
		FROM businesses
		WHERE id = $1
	`

	var b Business

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID,
		&b.FirebaseUID,
		&b.BusinessName,
		&b.OwnerName,
		&b.Phone,
		&b.Email,
		&b.Address,
		&b.LogoURL,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBusinessNotFound
		}

		return nil, fmt.Errorf("failed to get business: %w", err)
	}

	return &b, nil
}
