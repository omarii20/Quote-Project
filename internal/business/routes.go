package business

import "net/http"

func RegisterRoutes(handler *Handler) {
	http.HandleFunc("/businesses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateBusiness(w, r)

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
