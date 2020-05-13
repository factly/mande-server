package status

import "github.com/go-chi/chi"

// status request object
type status struct {
	Name string `json:"name"`
}

// Router - Group of product-status router
func Router() chi.Router {
	r := chi.NewRouter()
	r.Post("/", create)
	r.Get("/", list)

	return r
}
