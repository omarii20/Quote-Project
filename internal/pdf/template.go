package pdf

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/omarii20/Quote-Project/internal/business"
	"github.com/omarii20/Quote-Project/internal/quote"
	"github.com/shopspring/decimal"
)

//go:embed templates/quote.html
var templateFiles embed.FS

type ItemView struct {
	Description      string
	Quantity         string
	UnitPriceDisplay string
	TotalDisplay     string
}

type QuoteView struct {
	Business business.Business
	Quote    quote.Quote

	CustomerName  string
	CustomerPhone string
	CustomerEmail string

	CreatedDate string
	ValidUntil  string

	IsItems bool
	Items   []ItemView

	ManualSubtotalDisplay string
	ItemsSubtotalDisplay  string

	HasAdditionalAmount     bool
	AdditionalAmountDisplay string

	SubtotalDisplay string

	HasDiscount           bool
	DiscountAmountDisplay string

	VatRateDisplay   string
	VatAmountDisplay string
	TotalDisplay     string
}

func formatMoney(amount decimal.Decimal) string {
	return "₪" + amount.StringFixed(2)
}

func formatDate(value time.Time) string {
	return value.Format("02/01/2006")
}

func RenderQuoteHTML(
	businessData *business.Business,
	quoteData *quote.Quote,
) (string, error) {
	if businessData == nil || quoteData == nil {
		return "", fmt.Errorf("business and quote are required")
	}

	if businessData.ID != quoteData.BusinessID {
		return "", fmt.Errorf("quote does not belong to business")
	}

	view := QuoteView{
		Business: *businessData,
		Quote:    *quoteData,

		IsItems: quoteData.PricingMethod == "items",

		ItemsSubtotalDisplay: formatMoney(
			quoteData.ItemsSubtotal,
		),
		AdditionalAmountDisplay: formatMoney(
			quoteData.AdditionalAmount,
		),
		SubtotalDisplay: formatMoney(
			quoteData.Subtotal,
		),
		DiscountAmountDisplay: formatMoney(
			quoteData.DiscountAmount,
		),
		VatRateDisplay: quoteData.VATRate.StringFixed(2) + "%",
		VatAmountDisplay: formatMoney(
			quoteData.VATAmount,
		),
		TotalDisplay: formatMoney(
			quoteData.Total,
		),

		HasAdditionalAmount: !quoteData.AdditionalAmount.IsZero(),
		HasDiscount:         !quoteData.DiscountAmount.IsZero(),
	}

	if quoteData.Customer != nil {
		view.CustomerName = quoteData.Customer.Name
		view.CustomerPhone = quoteData.Customer.Phone

		view.CustomerEmail = quoteData.Customer.Email
	}

	if !quoteData.CreatedAt.IsZero() {
		view.CreatedDate = formatDate(quoteData.CreatedAt)
	}

	if quoteData.ValidUntil != nil {
		view.ValidUntil = formatDate(*quoteData.ValidUntil)
	}

	if view.IsItems {
		for _, item := range quoteData.Items {
			view.Items = append(view.Items, ItemView{
				Description: item.Description,
				Quantity:    item.Quantity.String(),
				UnitPriceDisplay: formatMoney(
					item.UnitPrice,
				),
				TotalDisplay: formatMoney(
					item.Total,
				),
			})
		}
	} else if quoteData.ManualSubtotal != nil {
		view.ManualSubtotalDisplay = formatMoney(
			*quoteData.ManualSubtotal,
		)
	}

	tmpl, err := template.ParseFS(
		templateFiles,
		"templates/quote.html",
	)
	if err != nil {
		return "", fmt.Errorf("failed to parse quote template: %w", err)
	}

	var output bytes.Buffer

	if err := tmpl.Execute(&output, view); err != nil {
		return "", fmt.Errorf("failed to render quote template: %w", err)
	}

	return strings.TrimSpace(output.String()), nil
}
