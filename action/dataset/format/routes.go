package format

import (
	"github.com/factly/data-portal-server/model"
	"github.com/go-chi/chi"
)

// datasetFormat request body
type datasetFormat struct {
	FormatID uint   `json:"format_id" validate:"required"`
	URL      string `json:"url" validate:"required"`
}

var userContext model.ContextKey = "dataset_user"

// AdminRouter - Group of tag router
func AdminRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)
	r.Post("/", create) // POST /format - create a new dataset format and persist it

	r.Route("/{format_id}", func(r chi.Router) {
		r.Delete("/", delete) // DELETE /format/{format_id} - delete a single dataset format by :format_id
	})

	return r
}
