package business

import (
	"encoding/json"
	"errors"
	"net/http"

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

// CreateBusiness handles POST /businesses.
func (h *Handler) CreateBusiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var b Business

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	firebaseUID, ok := auth.FirebaseUIDFromContext(r.Context())
	if !ok || firebaseUID == "" {
		http.Error(
			w,
			`{"error":"firebase uid not found in context"}`,
			http.StatusUnauthorized,
		)
		return
	}

	b.FirebaseUID = firebaseUID

	if err := h.service.CreateBusiness(r.Context(), &b); err != nil {
		if errors.Is(err, ErrBusinessAlreadyExists) {
			w.WriteHeader(http.StatusConflict)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "business already exists",
			})
			return
		}

		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"business": b,
	})
}

// GetBusiness handles GET /businesses/me.
func (h *Handler) GetBusiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "business id not found in context",
		})
		return
	}

	b, err := h.service.GetBusiness(r.Context(), businessID)
	if err != nil {
		if errors.Is(err, ErrBusinessNotFound) {
			w.WriteHeader(http.StatusNotFound)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "business not found",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "internal server error",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"business": b,
	})
}

// UpdateBusiness handles PATCH /businesses/me.
func (h *Handler) UpdateBusiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "business id not found in context",
		})
		return
	}

	var req UpdateBusinessRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	b, err := h.service.UpdateBusiness(r.Context(), businessID, &req)
	if err != nil {
		if errors.Is(err, ErrBusinessNotFound) {
			w.WriteHeader(http.StatusNotFound)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "business not found",
			})
			return
		}

		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"business": b,
	})
}
