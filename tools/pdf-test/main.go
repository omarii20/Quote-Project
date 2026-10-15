package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/omarii20/Quote-Project/internal/business"
	"github.com/omarii20/Quote-Project/internal/pdf"
	"github.com/omarii20/Quote-Project/internal/quote"
	"github.com/omarii20/Quote-Project/internal/quoteitem"
	"github.com/shopspring/decimal"
)

func strPtr(value string) *string {
	return &value
}

func money(value string) decimal.Decimal {
	return decimal.RequireFromString(value)
}

func main() {
	ctx := context.Background()

	generator := pdf.NewGenerator("/usr/bin/google-chrome")

	now := time.Now()
	validUntil := now.AddDate(0, 0, 30)

	businessData := &business.Business{
		ID:           1,
		BusinessName: "Omarii Software Solutions",
		OwnerName:    "מוחמד עומרי",
		Phone:        "050-1234567",
		Email:        strPtr("contact@example.com"),
		Address:      strPtr("נצרת, ישראל"),
	}

	customerData := &quote.QuoteCustomer{
		ID:    10,
		Name:  "חברת הדוגמה בע״מ",
		Phone: "052-1234567",
		Email: "customer@example.com",
	}

	// ----------------------------------------
	// Quote 1: Items pricing
	// ----------------------------------------

	itemsQuote := &quote.Quote{
		ID:            101,
		BusinessID:    businessData.ID,
		CustomerID:    customerData.ID,
		Customer:      customerData,
		QuoteNumber:   "Q-0001",
		Title:         strPtr("פיתוח אפליקציה"),
		Description:   strPtr("פיתוח אפליקציה לניהול לקוחות והצעות מחיר."),
		PricingMethod: "items",

		Items: []quoteitem.QuoteItem{
			{
				Description: "עיצוב ממשק משתמש",
				Quantity:    money("1"),
				UnitPrice:   money("1500"),
				Total:       money("1500"),
				Position:    1,
			},
			{
				Description: "פיתוח Backend ב-Go",
				Quantity:    money("2"),
				UnitPrice:   money("2000"),
				Total:       money("4000"),
				Position:    2,
			},
			{
				Description: "בדיקות והטמעה",
				Quantity:    money("1"),
				UnitPrice:   money("500"),
				Total:       money("500"),
				Position:    3,
			},
		},

		ItemsSubtotal:    money("6000"),
		AdditionalAmount: money("200"),
		Subtotal:         money("6200"),

		DiscountType:   strPtr("fixed"),
		DiscountValue:  money("200"),
		DiscountAmount: money("200"),

		VATRate:   money("18"),
		VATAmount: money("1080"),
		Total:     money("7080"),

		Status:     "draft",
		ValidUntil: &validUntil,
		Notes:      strPtr("המחיר כולל פיתוח, בדיקות והטמעה."),
		CreatedAt:  now,
	}

	// ----------------------------------------
	// Quote 2: Manual pricing
	// ----------------------------------------

	manualSubtotal := money("3500")

	manualQuote := &quote.Quote{
		ID:            102,
		BusinessID:    businessData.ID,
		CustomerID:    customerData.ID,
		Customer:      customerData,
		QuoteNumber:   "Q-0002",
		Title:         strPtr("תחזוקת מערכת"),
		Description:   strPtr("שירותי תחזוקה ותמיכה טכנית לתקופה של שלושה חודשים."),
		PricingMethod: "manual",

		ManualSubtotal:   &manualSubtotal,
		AdditionalAmount: money("0"),
		Subtotal:         money("3500"),

		DiscountType:   strPtr("fixed"),
		DiscountValue:  money("100"),
		DiscountAmount: money("100"),

		VATRate:   money("18"),
		VATAmount: money("612"),
		Total:     money("4012"),

		Status:     "draft",
		ValidUntil: &validUntil,
		Notes:      strPtr("התשלום בהתאם לתנאי ההתקשרות."),
		CreatedAt:  now,
	}

	tests := []struct {
		filename string
		data     *quote.Quote
	}{
		{
			filename: "quote-items.pdf",
			data:     itemsQuote,
		},
		{
			filename: "quote-manual.pdf",
			data:     manualQuote,
		},
	}

	for _, test := range tests {
		html, err := pdf.RenderQuoteHTML(
			businessData,
			test.data,
		)
		if err != nil {
			fmt.Printf("HTML rendering failed: %v\n", err)
			os.Exit(1)
		}

		pdfData, err := generator.Generate(ctx, html)
		if err != nil {
			fmt.Printf("PDF generation failed: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(
			test.filename,
			pdfData,
			0644,
		); err != nil {
			fmt.Printf("PDF save failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf(
			"Generated %s successfully (%d bytes)\n",
			test.filename,
			len(pdfData),
		)
	}
}
