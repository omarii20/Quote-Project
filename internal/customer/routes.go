package customer

import (
	"net/http"
	"strings"

	"github.com/omarii20/Quote-Project/internal/auth"
)

func RegisterRoutes(handler *Handler, authMiddleware *auth.Middleware, businessContext *auth.BusinessContext) {
	// POST /customers
	// GET  /customers
	http.Handle("/customers",
		authMiddleware.Authenticate(
			businessContext.Resolve(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {

					case http.MethodPost:
						handler.CreateCustomer(w, r)

					case http.MethodGet:
						handler.GetCustomersByBusinessID(w, r)

					default:
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusMethodNotAllowed)
					}
				}),
			),
		),
	)

	// GET    /customers/{id}
	// PATCH  /customers/{id}
	// DELETE /customers/{id}
	http.Handle("/customers/",
		authMiddleware.Authenticate(
			businessContext.Resolve(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					id := strings.TrimPrefix(r.URL.Path, "/customers/")

					switch r.Method {

					case http.MethodGet:
						handler.GetCustomer(w, r, id)

					case http.MethodPatch:
						handler.UpdateCustomer(w, r, id)

					case http.MethodDelete:
						handler.DeleteCustomer(w, r, id)

					default:
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusMethodNotAllowed)
					}
				}),
			),
		),
	)
}
