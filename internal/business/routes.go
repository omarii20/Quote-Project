package business

import (
	"net/http"
	"strings"
)

func RegisterRoutes(handler *Handler) {

	// /businesses
	http.HandleFunc("/businesses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateBusiness(w, r)

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// /businesses/{id}
	http.HandleFunc("/businesses/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:

			id := strings.TrimPrefix(r.URL.Path, "/businesses/")

			handler.GetBusiness(w, r, id)

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
