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

	r.Get("/", getTags)    // GET /tags - return list of tags
	r.Post("/", createTag) // POST /tags - create a new tag and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", getTagByID)   // GET /tags/{id} - read a single tag by :id
		r.Put("/", updateTag)    // PUT /tags/{id} - update a single tag by :id
		r.Delete("/", deleteTag) // DELETE /tags/{id} - delete a single tag by :id
	})

	return r
}
