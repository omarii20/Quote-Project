package quote

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/omarii20/Quote-Project/internal/quoteitem"
)

// Quote represents a quote in the system.
type Quote struct {
	ID         int64 `json:"id"`
	BusinessID int64 `json:"business_id"`
	CustomerID int64 `json:"customer_id"`

	QuoteNumber string  `json:"quote_number"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`

	PricingMethod string `json:"pricing_method"`

	ItemsSubtotal    decimal.Decimal  `json:"items_subtotal"`
	ManualSubtotal   *decimal.Decimal `json:"manual_subtotal,omitempty"`
	AdditionalAmount decimal.Decimal  `json:"additional_amount"`
	Subtotal         decimal.Decimal  `json:"subtotal"`

	DiscountType   *string         `json:"discount_type,omitempty"`
	DiscountValue  decimal.Decimal `json:"discount_value"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`

	VATRate   decimal.Decimal `json:"vat_rate"`
	VATAmount decimal.Decimal `json:"vat_amount"`

	Total decimal.Decimal `json:"total"`

	Status string `json:"status"`

	ValidUntil *time.Time `json:"valid_until,omitempty"`
	Notes      *string    `json:"notes,omitempty"`

	Items []quoteitem.QuoteItem `json:"items,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateQuoteRequest represents the request payload for updating a quote.
type UpdateQuoteRequest struct {
	CustomerID int64 `json:"customer_id"`

	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`

	PricingMethod string `json:"pricing_method"`

	ManualSubtotal   *decimal.Decimal `json:"manual_subtotal,omitempty"`
	AdditionalAmount decimal.Decimal  `json:"additional_amount"`

	DiscountType  *string         `json:"discount_type,omitempty"`
	DiscountValue decimal.Decimal `json:"discount_value"`

	VATRate decimal.Decimal `json:"vat_rate"`

	Status string `json:"status"`

	ValidUntil *time.Time `json:"valid_until,omitempty"`
	Notes      *string    `json:"notes,omitempty"`

	Items []quoteitem.QuoteItem `json:"items,omitempty"`
}

// UpdateQuoteStatusRequest represents the request payload
// for updating only the quote status.
type UpdateQuoteStatusRequest struct {
	Status string `json:"status"`
}
