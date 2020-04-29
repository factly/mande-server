package tag

import (
	"github.com/go-chi/chi"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /tags - return list of tags
	r.Post("/", create) // POST /tags - create a new tag and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", details)   // GET /tags/{id} - read a single tag by :id
		r.Delete("/", delete) // DELETE /tags/{id} - delete a single tag by :id
	})

	return r
}
