package quotedocument

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/omarii20/Quote-Project/internal/auth"
	"github.com/omarii20/Quote-Project/internal/quote"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetQuotePDF handles GET /quotes/{id}/pdf.
func (h *Handler) GetQuotePDF(
	w http.ResponseWriter,
	r *http.Request,
	id string,
) {
	// 1. Validate quote ID.
	quoteID, err := strconv.ParseInt(id, 10, 64)
	if err != nil || quoteID <= 0 {
		http.Error(w, "invalid quote id", http.StatusBadRequest)
		return
	}

	// 2. Get the authenticated business ID.
	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		http.Error(
			w,
			"business id not found in context",
			http.StatusInternalServerError,
		)
		return
	}

	// 3. Generate PDF.
	pdfData, quoteNumber, err := h.service.GenerateQuotePDF(
		r.Context(),
		quoteID,
		businessID,
	)
	if err != nil {
		if errors.Is(err, quote.ErrQuoteNotFound) {
			http.Error(w, "quote not found", http.StatusNotFound)
			return
		}

		http.Error(
			w,
			"failed to generate PDF",
			http.StatusInternalServerError,
		)
		return
	}

	// 4. Prepare a safe filename.
	safeNumber := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '-', r == '_':
			return r
		default:
			return -1
		}
	}, quoteNumber)

	if safeNumber == "" {
		safeNumber = strconv.FormatInt(quoteID, 10)
	}

	filename := fmt.Sprintf("quote-%s.pdf", safeNumber)

	// 5. Return PDF to the client.
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=%q", filename),
	)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfData)))

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pdfData)
}
