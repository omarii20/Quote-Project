package business

import (
	"net/http"

	"github.com/omarii20/Quote-Project/internal/auth"
)

func RegisterRoutes(
	handler *Handler,
	authMiddleware *auth.Middleware,
	businessContext *auth.BusinessContext,
) {

	// POST /businesses
	// User is authenticated, but does not have a business yet.
	http.Handle("/businesses",
		authMiddleware.Authenticate(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {

				case http.MethodPost:
					handler.CreateBusiness(w, r)

				default:
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			}),
		),
	)

	// GET   /businesses/me
	// PATCH /businesses/me
	// The business ID is resolved from the authenticated Firebase UID.
	http.Handle("/businesses/me",
		authMiddleware.Authenticate(
			businessContext.Resolve(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {

					case http.MethodGet:
						handler.GetBusiness(w, r)

					case http.MethodPatch:
						handler.UpdateBusiness(w, r)

					default:
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusMethodNotAllowed)
					}
				}),
			),
		),
	)
}
