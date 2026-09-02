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

// Update updates specific fields of an existing business.
func (r *Repository) Update(ctx context.Context, id int64, req *UpdateBusinessRequest) (*Business, error) {

	query := `
		UPDATE businesses
		SET
			business_name = COALESCE($1, business_name),
			owner_name = COALESCE($2, owner_name),
			phone = COALESCE($3, phone),
			email = COALESCE($4, email),
			address = COALESCE($5, address),
			logo_url = COALESCE($6, logo_url),
			updated_at = NOW()
		WHERE id = $7
		RETURNING
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
	`

	var b Business

	err := r.db.QueryRowContext(ctx, query, req.BusinessName, req.OwnerName, req.Phone, req.Email, req.Address, req.LogoURL, id).Scan(
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

		return nil, fmt.Errorf("failed to update business: %w", err)
	}

	return &b, nil
}

// GetBusinessIDByFirebaseUID retrieves a business ID by its Firebase UID.
func (r *Repository) GetBusinessIDByFirebaseUID(ctx context.Context, firebaseUID string) (int64, error) {

	query := `
		SELECT id
		FROM businesses
		WHERE firebase_uid = $1
	`

	var businessID int64

	err := r.db.QueryRowContext(
		ctx,
		query,
		firebaseUID,
	).Scan(&businessID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrBusinessNotFound
		}

		return 0, fmt.Errorf(
			"failed to get business id by firebase uid: %w",
			err,
		)
	}

	return businessID, nil
}
