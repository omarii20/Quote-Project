package customer

import (
	"net/http"
	"strings"
)

func RegisterRoutes(handler *Handler) {

	// POST /customers
	http.HandleFunc("/customers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:
			handler.CreateCustomer(w, r)

		case http.MethodGet:
			handler.GetCustomersByBusinessID(w, r)

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// GET   /customers/{id}
	// PATCH /customers/{id}
	http.HandleFunc("/customers/", func(w http.ResponseWriter, r *http.Request) {

		id := strings.TrimPrefix(r.URL.Path, "/customers/")

		switch r.Method {

		case http.MethodGet:
			handler.GetCustomer(w, r, id)

		case http.MethodPatch:
			handler.UpdateCustomer(w, r, id)

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
