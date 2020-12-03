package format

import (
	"github.com/factly/data-portal-server/model"
	"github.com/go-chi/chi"
)

// format request body
type format struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" `
	IsDefault   bool   `json:"is_default" `
}

var userContext model.ContextKey = "format_user"

// UserRouter - Group of format router
func UserRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list) // GET /formats - return list of formats

	r.Route("/{format_id}", func(r chi.Router) {
		r.Get("/", details) // GET /formats/{format_id} - read a single format by :format_id
	})
	return r
}

// AdminRouter - Group of format router
func AdminRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", create) // POST /formats - create a new format and persist it
	r.Get("/", list)    // GET /formats - return list of formats

	r.Route("/{format_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /formats/{format_id} - read a single format by :format_id
		r.Put("/", update)    // PUT /formats/{format_id} - update a single format by :format_id
		r.Delete("/", delete) // DELETE /formats/{format_id} - delete a single format by :format_id
	})
	return r
}
