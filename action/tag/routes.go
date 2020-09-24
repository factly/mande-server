package tag

import (
	"github.com/go-chi/chi"
)

// tag request body
type tag struct {
	Title string `json:"title" validate:"required"`
	Slug  string `json:"slug" validate:"required"`
}

// UserRouter - Group of tag router
func UserRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list) // GET /tags - return list of tags

	r.Route("/{tag_id}", func(r chi.Router) {
		r.Get("/", details) // GET /tags/{tag_id} - read a single tag by :tag_id
	})

	return r
}

// AdminRouter - Group of tag router
func AdminRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /tags - return list of tags
	r.Post("/", create) // POST /tags - create a new tag and persist it

	r.Route("/{tag_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /tags/{tag_id} - read a single tag by :tag_id
		r.Put("/", update)    // PUT /tags/{tag_id} - update a single tag by :tag_id
		r.Delete("/", delete) // DELETE /tags/{tag_id} - delete a single tag by :tag_id
	})

	return r
}
