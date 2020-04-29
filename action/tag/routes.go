package tag

import (
	"github.com/go-chi/chi"
)

// tag request body
type tag struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

// Router - Group of tag router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /tags - return list of tags
	r.Post("/", create) // POST /tags - create a new tag and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", detail)    // GET /tags/{id} - read a single tag by :id
		r.Put("/", update)    // PUT /tags/{id} - update a single tag by :id
		r.Delete("/", delete) // DELETE /tags/{id} - delete a single tag by :id
	})

	return r
}
