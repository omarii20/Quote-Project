package quote

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetNextQuoteNumber handles the HTTP request to retrieve the next quote number for a given business.
func (h *Handler) GetNextQuoteNumber(w http.ResponseWriter, r *http.Request) {
	businessIDStr := r.URL.Query().Get("business_id")

	businessID, err := strconv.ParseInt(businessIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid business id", http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

// CreateQuote handles the HTTP request to create a new quote.
func (h *Handler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var q Quote

	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateQuote(r.Context(), &q); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(q)
}

// GetQuote handles the HTTP request to retrieve a quote by its ID.
func (h *Handler) GetQuote(w http.ResponseWriter, r *http.Request, id string) {
	quoteID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid quote id", http.StatusBadRequest)
		return
	}

	q, err := h.service.GetQuoteByID(r.Context(), quoteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(q)
}

// GetQuotesByBusinessID handles the HTTP request to retrieve quotes for a specific business within a given period.
func (h *Handler) GetQuotesByBusinessID(w http.ResponseWriter, r *http.Request) {
	businessIDStr := r.URL.Query().Get("business_id")
	period := r.URL.Query().Get("period")

	businessID, err := strconv.ParseInt(businessIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid business id", http.StatusBadRequest)
		return
	}

	quotes, err := h.service.GetQuotesByBusinessID(r.Context(), businessID, period)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(quotes)
}

// DeleteQuote handles the HTTP request to delete a quote by its ID.
func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request, id string) {
	quoteID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid quote id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteQuote(r.Context(), quoteID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
