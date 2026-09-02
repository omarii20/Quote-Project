package quote

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/omarii20/Quote-Project/internal/auth"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetNextQuoteNumber handles GET /quotes/next-number.
func (h *Handler) GetNextQuoteNumber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		http.Error(
			w,
			"business id not found in context",
			http.StatusInternalServerError,
		)
		return
	}

	quoteNumber, err := h.service.GetNextQuoteNumber(r.Context(), businessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"quote_number": quoteNumber,
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

// CreateQuote handles POST /quotes.
func (h *Handler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var q Quote

	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		http.Error(
			w,
			"business id not found in context",
			http.StatusInternalServerError,
		)
		return
	}

	// Never trust business_id from the client.
	q.BusinessID = businessID

	if err := h.service.CreateQuote(r.Context(), &q); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(q)
}

// GetQuote handles GET /quotes/{id}.
func (h *Handler) GetQuote(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")

	quoteID, err := strconv.ParseInt(id, 10, 64)
	if err != nil || quoteID <= 0 {
		http.Error(w, "invalid quote id", http.StatusBadRequest)
		return
	}

	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		http.Error(
			w,
			"business id not found in context",
			http.StatusInternalServerError,
		)
		return
	}

	q, err := h.service.GetQuoteByID(r.Context(), quoteID, businessID)
	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			http.Error(w, "quote not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(q)
}

// GetQuotesByBusinessID handles GET /quotes?period={period}.
func (h *Handler) GetQuotesByBusinessID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		http.Error(
			w,
			"business id not found in context",
			http.StatusInternalServerError,
		)
		return
	}

	period := r.URL.Query().Get("period")

	quotes, err := h.service.GetQuotesByBusinessID(r.Context(), businessID, period)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(quotes)
}

// DeleteQuote handles DELETE /quotes/{id}.
func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request, id string) {
	quoteID, err := strconv.ParseInt(id, 10, 64)
	if err != nil || quoteID <= 0 {
		http.Error(w, "invalid quote id", http.StatusBadRequest)
		return
	}

	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		http.Error(
			w,
			"business id not found in context",
			http.StatusInternalServerError,
		)
		return
	}

	err = h.service.DeleteQuote(r.Context(), quoteID, businessID)
	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			http.Error(w, "quote not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateQuote handles PUT /quotes/{id}.
func (h *Handler) UpdateQuote(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")

	quoteID, err := strconv.ParseInt(id, 10, 64)
	if err != nil || quoteID <= 0 {
		http.Error(w, "invalid quote id", http.StatusBadRequest)
		return
	}

	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		http.Error(
			w,
			"business id not found in context",
			http.StatusInternalServerError,
		)
		return
	}

	var req UpdateQuoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedQuote, err := h.service.UpdateQuote(r.Context(), quoteID, businessID, &req)
	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			http.Error(w, "quote not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(updatedQuote)
}

// UpdateQuoteStatus handles PATCH /quotes/{id}/status.
func (h *Handler) UpdateQuoteStatus(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")

	quoteID, err := strconv.ParseInt(id, 10, 64)
	if err != nil || quoteID <= 0 {
		http.Error(w, "invalid quote id", http.StatusBadRequest)
		return
	}

	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		http.Error(
			w,
			"business id not found in context",
			http.StatusInternalServerError,
		)
		return
	}

	var req UpdateQuoteStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedQuote, err := h.service.UpdateQuoteStatus(r.Context(), quoteID, businessID, &req)
	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			http.Error(w, "quote not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(updatedQuote)
}
