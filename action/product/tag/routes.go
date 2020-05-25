package tag

import "github.com/go-chi/chi"

// ProductTag request body
type productTag struct {
	TagID uint `json:"tag_id"`
}

// Router - Group of product-tag router
func Router() chi.Router {
	r := chi.NewRouter()
	r.Post("/", create)

	r.Route("/{tag_id}", func(r chi.Router) {
		r.Delete("/", delete)
	})

	return r
}
