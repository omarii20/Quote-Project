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

// GetByBusinessID retrieves all customers associated with a specific business ID.
func (r *Repository) GetByBusinessID(ctx context.Context, businessID int64) ([]Customer, error) {

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
		WHERE business_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, businessID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}
	defer rows.Close()

	var customers []Customer

	for rows.Next() {
		var c Customer

		err := rows.Scan(
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
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}

		customers = append(customers, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read customers rows: %w", err)
	}

	return customers, nil
}

// BusinessExists checks if a business with the given ID exists in the database.
func (r *Repository) BusinessExists(ctx context.Context, businessID int64) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM businesses
			WHERE id = $1
		)
	`

	var exists bool

	err := r.db.QueryRowContext(
		ctx,
		query,
		businessID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check business existence: %w", err)
	}

	return exists, nil
}
