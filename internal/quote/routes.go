package quote

import (
	"net/http"
	"strings"

	"github.com/omarii20/Quote-Project/internal/auth"
)

func RegisterRoutes(
	handler *Handler,
	authMiddleware *auth.Middleware,
	businessContext *auth.BusinessContext,
) {

	// GET /quotes/next-number
	http.Handle("/quotes/next-number",
		authMiddleware.Authenticate(
			businessContext.Resolve(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case http.MethodGet:
						handler.GetNextQuoteNumber(w, r)

					default:
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusMethodNotAllowed)
					}
				}),
			),
		),
	)

	// POST /quotes
	// GET /quotes
	http.Handle("/quotes",
		authMiddleware.Authenticate(
			businessContext.Resolve(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case http.MethodPost:
						handler.CreateQuote(w, r)

					case http.MethodGet:
						handler.GetQuotesByBusinessID(w, r)

					default:
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusMethodNotAllowed)
					}
				}),
			),
		),
	)

	// GET /quotes/{id}
	// PUT /quotes/{id}
	// DELETE /quotes/{id}
	// PATCH /quotes/{id}/status
	http.Handle("/quotes/",
		authMiddleware.Authenticate(
			businessContext.Resolve(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					path := strings.TrimPrefix(r.URL.Path, "/quotes/")

					if strings.HasSuffix(path, "/status") {
						id := strings.TrimSuffix(path, "/status")

						switch r.Method {
						case http.MethodPatch:
							handler.UpdateQuoteStatus(w, r, id)

						default:
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusMethodNotAllowed)
						}

						return
					}

					id := path

					switch r.Method {
					case http.MethodGet:
						handler.GetQuote(w, r, id)

					case http.MethodPut:
						handler.UpdateQuote(w, r, id)

					case http.MethodDelete:
						handler.DeleteQuote(w, r, id)

					default:
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusMethodNotAllowed)
					}
				}),
			),
		),
	)
}
