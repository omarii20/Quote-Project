package business

import (
	"net/http"
	"strings"
)

func RegisterRoutes(handler *Handler) {

	// POST/businesses
	http.HandleFunc("/businesses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateBusiness(w, r)

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// GET-PATCH /businesses/{id}
	http.HandleFunc("/businesses/", func(w http.ResponseWriter, r *http.Request) {

		id := strings.TrimPrefix(r.URL.Path, "/businesses/")

		switch r.Method {

		case http.MethodGet:
			handler.GetBusiness(w, r, id)

		case http.MethodPatch:
			handler.UpdateBusiness(w, r, id)

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

}
