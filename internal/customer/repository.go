package customer

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

// Create inserts a new customer into the database.
func (r *Repository) Create(ctx context.Context, c *Customer) error {
	query := `
		INSERT INTO customers (
			business_id,
			name,
			phone,
			email,
			address,
			notes
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query, c.BusinessID, c.Name, c.Phone, c.Email, c.Address, c.Notes).Scan(
		&c.ID,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}

	return nil
}

// GetByID retrieves a customer by its ID.
func (r *Repository) GetByID(ctx context.Context, id int64) (*Customer, error) {
	query := `
		SELECT
			id,
			business_id,
			name,
			phone,
			email,
			address,
			notes,
			created_at,
			updated_at
		FROM customers
		WHERE id = $1
	`

	var c Customer

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.BusinessID,
		&c.Name,
		&c.Phone,
		&c.Email,
		&c.Address,
		&c.Notes,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCustomerNotFound
		}

		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &c, nil
}
