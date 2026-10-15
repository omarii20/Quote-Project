package pdf

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type Generator struct {
	chromePath string
}

func NewGenerator(chromePath string) *Generator {
	return &Generator{
		chromePath: chromePath,
	}
}

func (g *Generator) Generate(
	ctx context.Context,
	htmlContent string,
) ([]byte, error) {

	if strings.TrimSpace(htmlContent) == "" {
		return nil, fmt.Errorf("HTML content is empty")
	}

	options := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(g.chromePath),
		chromedp.Flag("headless", true),
	)

	allocatorCtx, cancelAllocator := chromedp.NewExecAllocator(
		ctx,
		options...,
	)
	defer cancelAllocator()

	browserCtx, cancelBrowser := chromedp.NewContext(allocatorCtx)
	defer cancelBrowser()

	timeoutCtx, cancelTimeout := context.WithTimeout(
		browserCtx,
		30*time.Second,
	)
	defer cancelTimeout()

	// Load the HTML into Chrome.
	htmlURL := "data:text/html;charset=utf-8," +
		url.PathEscape(htmlContent)

	if err := chromedp.Do(
		timeoutCtx,
		chromedp.Navigate(htmlURL),
	); err != nil {
		return nil, fmt.Errorf("failed to load HTML: %w", err)
	}

	// Wait for web fonts to finish loading.
	_, err := chromedp.Run(
		timeoutCtx,
		chromedp.Evaluate[bool](
			`document.fonts.status === "loaded"`,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to check fonts: %w", err)
	}
	// Convert the rendered page to PDF.
	pdfData, err := chromedp.Run(
		timeoutCtx,
		chromedp.PrintToPDF(
			chromedp.PDFPaper(chromedp.PaperA4),
			chromedp.PDFPrintBackground(),
			chromedp.PDFPreferCSSPageSize(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfData, nil
}
