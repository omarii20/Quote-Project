package quotedocument

import (
	"context"
	"fmt"

	"github.com/omarii20/Quote-Project/internal/business"
	"github.com/omarii20/Quote-Project/internal/pdf"
	"github.com/omarii20/Quote-Project/internal/quote"
)

// Service coordinates the generation of quote PDF documents.
type Service struct {
	quoteService *quote.Service
	businessRepo *business.Repository
	pdfGenerator *pdf.Generator
}

func NewService(
	quoteService *quote.Service,
	businessRepo *business.Repository,
	pdfGenerator *pdf.Generator,
) *Service {
	return &Service{
		quoteService: quoteService,
		businessRepo: businessRepo,
		pdfGenerator: pdfGenerator,
	}
}

// GenerateQuotePDF generates a PDF for a quote
// belonging to the authenticated business.
func (s *Service) GenerateQuotePDF(
	ctx context.Context,
	quoteID int64,
	businessID int64,
) ([]byte, string, error) {

	if quoteID <= 0 || businessID <= 0 {
		return nil, "", fmt.Errorf("invalid quote or business id")
	}

	// 1. Retrieve the quote with business-scoped authorization.
	quoteData, err := s.quoteService.GetQuoteByID(
		ctx,
		quoteID,
		businessID,
	)
	if err != nil {
		return nil, "", err
	}

	// 2. Retrieve the authenticated business details.
	businessData, err := s.businessRepo.GetByID(
		ctx,
		businessID,
	)
	if err != nil {
		return nil, "", fmt.Errorf(
			"failed to retrieve business: %w",
			err,
		)
	}

	// 3. Render the quote HTML template.
	htmlContent, err := pdf.RenderQuoteHTML(
		businessData,
		quoteData,
	)
	if err != nil {
		return nil, "", fmt.Errorf(
			"failed to render quote HTML: %w",
			err,
		)
	}

	// 4. Generate PDF using Google Chrome.
	pdfData, err := s.pdfGenerator.Generate(
		ctx,
		htmlContent,
	)
	if err != nil {
		return nil, "", fmt.Errorf(
			"failed to generate quote PDF: %w",
			err,
		)
	}

	return pdfData, quoteData.QuoteNumber, nil
}
