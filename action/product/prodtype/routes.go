package prodtype

import "github.com/go-chi/chi"

// productType request body
type productType struct {
	Name string `json:"name"`
}

// Router - Group of product category router
func Router() chi.Router {
	r := chi.NewRouter()
	r.Post("/", create)
	r.Get("/", list)

	return r
}
