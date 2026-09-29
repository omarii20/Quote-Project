package quoteitem

import (
	"time"

	"github.com/shopspring/decimal"
)

type QuoteItem struct {
	ID      int64 `json:"id"`
	QuoteID int64 `json:"quote_id"`

	Description string          `json:"description"`
	Quantity    decimal.Decimal `json:"quantity"`
	UnitPrice   decimal.Decimal `json:"unit_price"`

	Total           decimal.Decimal `json:"total"`
	TotalOverridden bool            `json:"total_overridden"`

	Position int `json:"position"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
