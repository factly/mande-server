package catalog

import (
	"time"

	"github.com/factly/mande-server/model"
	"github.com/go-chi/chi"
)

// Catalog request body
type catalog struct {
	Title            string     `json:"title" validate:"required"`
	Slug             string     `json:"slug" validate:"required"`
	Price            int        `json:"price" validate:"required"`
	Description      string     `json:"description"`
	CurrencyID       uint       `json:"currency_id"`
	FeaturedMediumID uint       `json:"featured_medium_id"`
	PublishedDate    *time.Time `json:"published_date"`
	ProductIDs       []uint     `json:"product_ids"`
}

var userContext model.ContextKey = "catalog_user"

// PublicRouter - Group of catalog user router
func PublicRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list) // GET /catalogs - return list of catalogs
	r.Get("/my", my) // GET /catalogs/my - return list of catalogs owned by user
	r.Route("/{catalog_id}", func(r chi.Router) {
		r.Get("/", details) // GET /catalogs/{catalog_id} - read a single catalog by :catalog_id
	})

	return r
}

// AdminRouter - Group of catalog user router
func AdminRouter() chi.Router {
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
