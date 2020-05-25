package category

import "github.com/go-chi/chi"

// ProductCategory request body
type productCategory struct {
	CategoryID uint `json:"category_id"`
}

// Router - Group of product-category router
func Router() chi.Router {
	r := chi.NewRouter()
	r.Post("/", create)

	r.Route("/{category_id}", func(r chi.Router) {
		r.Delete("/", delete)
	})

	return r
}
