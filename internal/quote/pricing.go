package quote

import (
	"errors"
	"strings"

	"github.com/shopspring/decimal"
)

// calculateItemsPricing calculates the pricing for a quote using the "items" pricing method.
func calculateItemsPricing(q *Quote) error {
	itemsSubtotal := decimal.Zero

	for i := range q.Items {
		item := &q.Items[i]

		if strings.TrimSpace(item.Description) == "" {
			return errors.New("item description is required")
		}

		if item.Quantity.LessThanOrEqual(decimal.Zero) {
			return errors.New("item quantity must be greater than zero")
		}

		if item.UnitPrice.IsNegative() {
			return errors.New("item unit price cannot be negative")
		}

		if !item.TotalOverridden {
			item.Total = item.Quantity.Mul(item.UnitPrice)
		}

		if item.Total.IsNegative() {
			return errors.New("item total cannot be negative")
		}

		itemsSubtotal = itemsSubtotal.Add(item.Total)

		if item.Position == 0 {
			item.Position = i + 1
		}
	}

	q.ItemsSubtotal = itemsSubtotal

	if q.AdditionalAmount.IsNegative() {
		return errors.New("additional amount cannot be negative")
	}

	q.Subtotal = q.ItemsSubtotal.Add(q.AdditionalAmount)

	return nil
}

// calculateManualPricing calculates the pricing for a quote using the "manual" pricing method.
func calculateManualPricing(q *Quote) error {
	if q.ManualSubtotal == nil {
		return errors.New("manual subtotal is required")
	}

	if q.ManualSubtotal.IsNegative() {
		return errors.New("manual subtotal cannot be negative")
	}

	q.ItemsSubtotal = decimal.Zero
	q.AdditionalAmount = decimal.Zero
	q.Subtotal = *q.ManualSubtotal

	return nil
}

// calculateFinalPrice calculates the final price for a quote after applying discounts and VAT.
func calculateFinalPrice(q *Quote) error {
	if q.DiscountValue.IsNegative() {
		return errors.New("discount value cannot be negative")
	}

	if q.VATRate.IsNegative() {
		return errors.New("vat rate cannot be negative")
	}

	q.DiscountAmount = decimal.Zero

	if q.DiscountType != nil {
		switch *q.DiscountType {
		case "percent":
			if q.DiscountValue.GreaterThan(decimal.NewFromInt(100)) {
				return errors.New("discount percent cannot be greater than 100")
			}

			q.DiscountAmount = q.Subtotal.
				Mul(q.DiscountValue).
				Div(decimal.NewFromInt(100))

		case "fixed":
			q.DiscountAmount = q.DiscountValue

		default:
			return errors.New("invalid discount type")
		}
	}

	if q.DiscountAmount.GreaterThan(q.Subtotal) {
		return errors.New("discount cannot be greater than subtotal")
	}

	afterDiscount := q.Subtotal.Sub(q.DiscountAmount)

	q.VATAmount = afterDiscount.
		Mul(q.VATRate).
		Div(decimal.NewFromInt(100))

	q.Total = afterDiscount.Add(q.VATAmount)

	return nil
}
