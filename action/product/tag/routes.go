package tag

import "github.com/go-chi/chi"

// ProductTag request body
type productTags struct {
	TagID uint `json:"tag_id"`
}

// Router - Group of product category router
func Router() chi.Router {
	r := chi.NewRouter()
	r.Post("/", create)

	r.Route("/{cid}", func(r chi.Router) {
		r.Delete("/", delete)
	})

	return r
}
