package tag

import (
	"github.com/go-chi/chi"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", getTags)    // GET /tag - return list of tags
	r.Post("/", createTag) // POST /tag - create a new tag and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", getTagByID)   // GET /tags/{id} - read a single tag by :id
		r.Delete("/", deleteTag) // DELETE /tags/{id} - delete a single tag by :id
	})

	return r
}
