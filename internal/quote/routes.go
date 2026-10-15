package quote

import (
	"net/http"
	"strings"

	"github.com/omarii20/Quote-Project/internal/auth"
)

// PDFHandler handles GET /quotes/{id}/pdf.
type PDFHandler func(
	w http.ResponseWriter,
	r *http.Request,
	id string,
)

func RegisterRoutes(
	handler *Handler,
	authMiddleware *auth.Middleware,
	businessContext *auth.BusinessContext,
	pdfHandler PDFHandler,
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
	// GET /quotes/{id}/pdf
	http.Handle("/quotes/",
		authMiddleware.Authenticate(
			businessContext.Resolve(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					path := strings.TrimPrefix(r.URL.Path, "/quotes/")

					// GET /quotes/{id}/pdf
					if strings.HasSuffix(path, "/pdf") {
						id := strings.TrimSuffix(path, "/pdf")

						if r.Method != http.MethodGet {
							w.WriteHeader(http.StatusMethodNotAllowed)
							return
						}

						pdfHandler(w, r, id)
						return
					}

					// PATCH /quotes/{id}/status
					if strings.HasSuffix(path, "/status") {
						id := strings.TrimSuffix(path, "/status")

						if r.Method != http.MethodPatch {
							w.WriteHeader(http.StatusMethodNotAllowed)
							return
						}

						handler.UpdateQuoteStatus(w, r, id)
						return
					}

					// /quotes/{id}
					id := path

					switch r.Method {
					case http.MethodGet:
						handler.GetQuote(w, r, id)

					case http.MethodPut:
						handler.UpdateQuote(w, r, id)

					case http.MethodDelete:
						handler.DeleteQuote(w, r, id)

					default:
						w.WriteHeader(http.StatusMethodNotAllowed)
					}
				}),
			),
		),
	)
}
