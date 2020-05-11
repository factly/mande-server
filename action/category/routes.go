package category

import "github.com/go-chi/chi"

// Category request body
type category struct {
	Title    string `json:"title"`
	Slug     string `json:"slug"`
	ParentID uint   `json:"parent_id"`
}

// Router - Group of currency router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /categories - return list of categories
	r.Post("/", create) // POST /categories - create a new category and persist it

	r.Route("/{category_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /categories/{category_id} - read a single category by :category_id
		r.Put("/", update)    // PUT /categories/{category_id} - update a single category by :category_id
		r.Delete("/", delete) // DELETE /categories/{category_id} - delete a single category by :category_id
	})

	return r
}
