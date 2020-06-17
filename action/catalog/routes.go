package catalog

import (
	"time"

	"github.com/factly/data-portal-server/model"
	"github.com/go-chi/chi"
)

// Catalog request body
type catalog struct {
	Title           string    `json:"title" validate:"required"`
	Description     string    `json:"description" `
	Price           int       `json:"price" validate:"required"`
	FeaturedMediaID uint      `json:"featured_media_id"`
	PublishedDate   time.Time `json:"published_date" validate:"required"`
	ProductIDs      []uint    `json:"product_ids"`
}

type catalogData struct {
	model.Catalog
	Products []model.Product `json:"products"`
}

// Router - Group of catalog router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /catalogs - return list of catalogs
	r.Post("/", create) // POST /catalogs - create a new catalog and persist it

	r.Route("/{catalog_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /catalogs/{catalog_id} - read a single catalog by :catalog_id
		r.Put("/", update)    // PUT /catalogs/{catalog_id} - update a single catalog by :catalog_id
		r.Delete("/", delete) // DELETE /catalogs/{catalog_id} - delete a single catalog by :catalog_id
	})

	return r
}
