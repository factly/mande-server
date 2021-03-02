package medium

import (
	"github.com/factly/data-portal-server/model"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// medium request body
type medium struct {
	Name        string         `json:"name" validate:"required"`
	Slug        string         `json:"slug"`
	Type        string         `json:"type"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Caption     string         `json:"caption"`
	AltText     string         `json:"alt_text"`
	FileSize    int            `json:"file_size" validate:"required"`
	URL         postgres.Jsonb `json:"url" swaggertype:"primitive,string"`
	Dimensions  string         `json:"dimensions"`
}

var userContext model.ContextKey = "medium_user"

// UserRouter - Group of medium router
func UserRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list) // GET /media - return list of media

	r.Route("/{medium_id}", func(r chi.Router) {
		r.Get("/", details) // GET /media/{medium_id} - read a single medium by :medium_id
	})

	return r
}

// AdminRouter - Group of medium router
func AdminRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /media - return list of media
	r.Post("/", create) // POST /media - create a new medium and persist it

	r.Route("/{medium_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /media/{medium_id} - read a single medium by :medium_id
		r.Put("/", update)    // PUT /media/{medium_id} - update a single medium by :medium_id
		r.Delete("/", delete) // DELETE /media/{medium_id} - delete a single medium by :medium_id
	})

	return r
}
