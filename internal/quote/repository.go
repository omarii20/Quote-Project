package quote

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/omarii20/Quote-Project/internal/quoteitem"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// CustomerBelongsToBusiness checks if a customer belongs to a specific business.
func (r *Repository) CustomerBelongsToBusiness(ctx context.Context, customerID int64, businessID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM customers
			WHERE id = $1
			  AND business_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, customerID, businessID).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check customer business: %w", err)
	}

	return exists, nil
}

// BusinessExists checks if a business exists in the database.
func (r *Repository) BusinessExists(ctx context.Context, businessID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM businesses
			WHERE id = $1
		)
	`

	var exists bool

	err := r.db.QueryRowContext(ctx, query, businessID).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check business: %w", err)
	}

	return exists, nil
}

// Create inserts a new quote and its associated items into the database.
func (r *Repository) Create(ctx context.Context, q *Quote) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure the transaction is rolled back in case of an error
	defer tx.Rollback()

	// Lock the business row to prevent concurrent quote number generation
	var lockedBusinessID int64
	err = tx.QueryRowContext(
		ctx,
		`
			SELECT id
			FROM businesses
			WHERE id = $1
			FOR UPDATE
		`,
		q.BusinessID,
	).Scan(&lockedBusinessID)

	if err != nil {
		return fmt.Errorf("failed to lock business: %w", err)
	}

	// Generate the next quote number for the business
	var nextNumber int

	err = tx.QueryRowContext(
		ctx,
		`
			SELECT COALESCE(
				MAX(
					CAST(
						SUBSTRING(quote_number FROM 3)
						AS INTEGER
					)
				),
				0
			) + 1
			FROM quotes
			WHERE business_id = $1
		`,
		q.BusinessID,
	).Scan(&nextNumber)

	if err != nil {
		return fmt.Errorf("failed to get next quote number: %w", err)
	}

	q.QuoteNumber = fmt.Sprintf("Q-%04d", nextNumber)

	quoteQuery := `
		INSERT INTO quotes (
			business_id,
			customer_id,
			quote_number,
			title,
			description,
			pricing_method,
			items_subtotal,
			manual_subtotal,
			additional_amount,
			subtotal,
			discount_type,
			discount_value,
			discount_amount,
			vat_rate,
			vat_amount,
			total,
			status,
			valid_until,
			notes
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15,
			$16, $17, $18, $19
		)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRowContext(
		ctx,
		quoteQuery,
		q.BusinessID,
		q.CustomerID,
		q.QuoteNumber,
		q.Title,
		q.Description,
		q.PricingMethod,
		q.ItemsSubtotal,
		q.ManualSubtotal,
		q.AdditionalAmount,
		q.Subtotal,
		q.DiscountType,
		q.DiscountValue,
		q.DiscountAmount,
		q.VATRate,
		q.VATAmount,
		q.Total,
		q.Status,
		q.ValidUntil,
		q.Notes,
	).Scan(
		&q.ID,
		&q.CreatedAt,
		&q.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create quote: %w", err)
	}

	itemQuery := `
		INSERT INTO quote_items (
			quote_id,
			description,
			quantity,
			unit_price,
			total,
			total_overridden,
			position
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	for i := range q.Items {
		item := &q.Items[i]

		item.QuoteID = q.ID

		err = tx.QueryRowContext(
			ctx,
			itemQuery,
			item.QuoteID,
			item.Description,
			item.Quantity,
			item.UnitPrice,
			item.Total,
			item.TotalOverridden,
			item.Position,
		).Scan(
			&item.ID,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to create quote item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetNextQuoteNumber retrieves the next available quote number for a given business.
func (r *Repository) GetNextQuoteNumber(ctx context.Context, businessID int64) (string, error) {
	var nextNumber int

	query := `
		SELECT COALESCE(
			MAX(
				CAST(
					SUBSTRING(quote_number FROM 3)
					AS INTEGER
				)
			),
			0
		) + 1
		FROM quotes
		WHERE business_id = $1
	`

	err := r.db.QueryRowContext(ctx, query, businessID).Scan(&nextNumber)
	if err != nil {
		return "", fmt.Errorf("failed to get next quote number: %w", err)
	}

	return fmt.Sprintf("Q-%04d", nextNumber), nil
}

// GetByID retrieves a quote by its ID from the database.
func (r *Repository) GetByID(ctx context.Context, quoteID int64) (*Quote, error) {
	query := `
		SELECT
			id,
			business_id,
			customer_id,
			quote_number,
			title,
			description,
			pricing_method,
			items_subtotal,
			manual_subtotal,
			additional_amount,
			subtotal,
			discount_type,
			discount_value,
			discount_amount,
			vat_rate,
			vat_amount,
			total,
			status,
			valid_until,
			notes,
			created_at,
			updated_at
		FROM quotes
		WHERE id = $1
	`

	var q Quote

	err := r.db.QueryRowContext(ctx, query, quoteID).Scan(
		&q.ID,
		&q.BusinessID,
		&q.CustomerID,
		&q.QuoteNumber,
		&q.Title,
		&q.Description,
		&q.PricingMethod,
		&q.ItemsSubtotal,
		&q.ManualSubtotal,
		&q.AdditionalAmount,
		&q.Subtotal,
		&q.DiscountType,
		&q.DiscountValue,
		&q.DiscountAmount,
		&q.VATRate,
		&q.VATAmount,
		&q.Total,
		&q.Status,
		&q.ValidUntil,
		&q.Notes,
		&q.CreatedAt,
		&q.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("quote not found")
		}

		return nil, fmt.Errorf("failed to get quote: %w", err)
	}

	// Retrieve associated quote items
	itemsQuery := `
		SELECT
			id,
			quote_id,
			description,
			quantity,
			unit_price,
			total,
			total_overridden,
			position,
			created_at,
			updated_at
		FROM quote_items
		WHERE quote_id = $1
		ORDER BY position ASC
	`

	rows, err := r.db.QueryContext(ctx, itemsQuery, quoteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item quoteitem.QuoteItem

		err := rows.Scan(
			&item.ID,
			&item.QuoteID,
			&item.Description,
			&item.Quantity,
			&item.UnitPrice,
			&item.Total,
			&item.TotalOverridden,
			&item.Position,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quote item: %w", err)
		}

		q.Items = append(q.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while reading quote items: %w", err)
	}

	return &q, nil
}

// GetByBusinessID retrieves (day, week, month, or all) quotes for a specific business within a given period.
func (r *Repository) GetByBusinessID(ctx context.Context, businessID int64, period string) ([]Quote, error) {

	query := `
		SELECT
			id,
			business_id,
			customer_id,
			quote_number,
			title,
			description,
			pricing_method,
			items_subtotal,
			manual_subtotal,
			additional_amount,
			subtotal,
			discount_type,
			discount_value,
			discount_amount,
			vat_rate,
			vat_amount,
			total,
			status,
			valid_until,
			notes,
			created_at,
			updated_at
		FROM quotes
		WHERE business_id = $1
	`

	switch period {
	case "today":
		query += ` AND created_at >= CURRENT_DATE`

	case "week":
		query += ` AND created_at >= NOW() - INTERVAL '7 days'`

	case "month":
		query += ` AND created_at >= NOW() - INTERVAL '30 days'`

	case "all":
		// No date filter

	default:
		return nil, errors.New("invalid period")
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, businessID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes: %w", err)
	}
	defer rows.Close()

	var quotes []Quote

	for rows.Next() {
		var q Quote

		err := rows.Scan(
			&q.ID,
			&q.BusinessID,
			&q.CustomerID,
			&q.QuoteNumber,
			&q.Title,
			&q.Description,
			&q.PricingMethod,
			&q.ItemsSubtotal,
			&q.ManualSubtotal,
			&q.AdditionalAmount,
			&q.Subtotal,
			&q.DiscountType,
			&q.DiscountValue,
			&q.DiscountAmount,
			&q.VATRate,
			&q.VATAmount,
			&q.Total,
			&q.Status,
			&q.ValidUntil,
			&q.Notes,
			&q.CreatedAt,
			&q.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quote: %w", err)
		}

		quotes = append(quotes, q)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while reading quotes: %w", err)
	}

	return quotes, nil
}

// Delete removes a quote from the database by its ID.
func (r *Repository) Delete(ctx context.Context, quoteID int64) error {
	result, err := r.db.ExecContext(
		ctx,
		`
			DELETE FROM quotes
			WHERE id = $1
		`,
		quoteID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete quote: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check deleted quote: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("quote not found")
	}

	return nil
}

// UpdateQuote updates an existing quote and its associated items in the database.
func (r *Repository) UpdateQuote(ctx context.Context, q *Quote) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	quoteQuery := `
		UPDATE quotes
		SET
			customer_id = $1,
			title = $2,
			description = $3,
			pricing_method = $4,
			items_subtotal = $5,
			manual_subtotal = $6,
			additional_amount = $7,
			subtotal = $8,
			discount_type = $9,
			discount_value = $10,
			discount_amount = $11,
			vat_rate = $12,
			vat_amount = $13,
			total = $14,
			status = $15,
			valid_until = $16,
			notes = $17,
			updated_at = NOW()
		WHERE id = $18
		RETURNING updated_at
	`

	err = tx.QueryRowContext(
		ctx,
		quoteQuery,
		q.CustomerID,
		q.Title,
		q.Description,
		q.PricingMethod,
		q.ItemsSubtotal,
		q.ManualSubtotal,
		q.AdditionalAmount,
		q.Subtotal,
		q.DiscountType,
		q.DiscountValue,
		q.DiscountAmount,
		q.VATRate,
		q.VATAmount,
		q.Total,
		q.Status,
		q.ValidUntil,
		q.Notes,
		q.ID,
	).Scan(
		&q.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("quote not found")
		}

		return fmt.Errorf("failed to update quote: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`
			DELETE FROM quote_items
			WHERE quote_id = $1
		`,
		q.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete old quote items: %w", err)
	}

	itemQuery := `
		INSERT INTO quote_items (
			quote_id,
			description,
			quantity,
			unit_price,
			total,
			total_overridden,
			position
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	for i := range q.Items {
		item := &q.Items[i]

		item.QuoteID = q.ID

		err = tx.QueryRowContext(
			ctx,
			itemQuery,
			item.QuoteID,
			item.Description,
			item.Quantity,
			item.UnitPrice,
			item.Total,
			item.TotalOverridden,
			item.Position,
		).Scan(
			&item.ID,
			&item.CreatedAt,
			&item.UpdatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to create updated quote item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit quote update: %w", err)
	}

	return nil
}

// UpdateQuoteStatus updates the status of an existing quote in the database.
func (r *Repository) UpdateQuoteStatus(ctx context.Context, quoteID int64, status string) error {
	result, err := r.db.ExecContext(
		ctx,
		`
			UPDATE quotes
			SET
				status = $1,
				updated_at = NOW()
			WHERE id = $2
		`,
		status,
		quoteID,
	)
	if err != nil {
		return fmt.Errorf("failed to update quote status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check updated quote status: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("quote not found")
	}

	return nil
}
