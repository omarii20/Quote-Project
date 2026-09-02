package customer

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

// CreateCustomer handles POST /customers.
func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var c Customer

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "business id not found in context",
		})
		return
	}

	c.BusinessID = businessID

	if err := h.service.CreateCustomer(r.Context(), &c); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"customer": c,
	})
}

// GetCustomer handles GET /customers/{id}.
func (h *Handler) GetCustomer(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")

	customerID, err := strconv.ParseInt(id, 10, 64)
	if err != nil || customerID <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid customer id",
		})
		return
	}

	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "business id not found in context",
		})
		return
	}

	c, err := h.service.GetCustomer(
		r.Context(),
		customerID,
		businessID,
	)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			w.WriteHeader(http.StatusNotFound)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "customer not found",
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
		"customer": c,
	})
}

// GetCustomersByBusinessID handles GET /customers.
func (h *Handler) GetCustomersByBusinessID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "business id not found in context",
		})
		return
	}

	customers, err := h.service.GetCustomersByBusinessID(
		r.Context(),
		businessID,
	)
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
		"success":   true,
		"customers": customers,
	})
}

// UpdateCustomer handles PATCH /customers/{id}.
func (h *Handler) UpdateCustomer(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")

	customerID, err := strconv.ParseInt(id, 10, 64)
	if err != nil || customerID <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid customer id",
		})
		return
	}

	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "business id not found in context",
		})
		return
	}

	var req UpdateCustomerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	c, err := h.service.UpdateCustomer(
		r.Context(),
		customerID,
		businessID,
		&req,
	)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			w.WriteHeader(http.StatusNotFound)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "customer not found",
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
		"customer": c,
	})
}

// DeleteCustomer handles DELETE /customers/{id}.
func (h *Handler) DeleteCustomer(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")

	customerID, err := strconv.ParseInt(id, 10, 64)
	if err != nil || customerID <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid customer id",
		})
		return
	}

	businessID, ok := auth.BusinessIDFromContext(r.Context())
	if !ok || businessID <= 0 {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "business id not found in context",
		})
		return
	}

	err = h.service.DeleteCustomer(
		r.Context(),
		customerID,
		businessID,
	)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			w.WriteHeader(http.StatusNotFound)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "customer not found",
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
		"success": true,
		"message": "customer deleted successfully",
	})
}
