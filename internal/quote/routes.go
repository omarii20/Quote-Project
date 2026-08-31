package quote

import (
	"net/http"
	"strings"
)

func RegisterRoutes(handler *Handler) {
	// Register the route for retrieving the next quote number
	http.HandleFunc("/quotes/next-number", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetNextQuoteNumber(w, r)

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Register the route for creating a new quote
	http.HandleFunc("/quotes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateQuote(w, r)

		case http.MethodGet:
			handler.GetQuotesByBusinessID(w, r)

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Register the route for retrieving a quote by its ID
	http.HandleFunc("/quotes/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/quotes/")

		switch r.Method {
		case http.MethodGet:
			handler.GetQuote(w, r, id)

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
