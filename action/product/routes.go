package product

import (
	"github.com/factly/data-portal-api/action/product/category"
	"github.com/factly/data-portal-api/action/product/prodtype"
	"github.com/factly/data-portal-api/action/product/status"
	"github.com/factly/data-portal-api/action/product/tag"
	"github.com/go-chi/chi"
)

type product struct {
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Price         int    `json:"price"`
	ProductTypeID uint   `json:"product_type_id"`
	StatusID      uint   `json:"status_id"`
	CurrencyID    uint   `json:"currency_id"`
}

// Router - Group of product router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/", create)
	r.Get("/", list)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", detail)
		r.Delete("/", delete)
		r.Put("/", update)
		r.Mount("/type", prodtype.Router())     // product-type router
		r.Mount("/status", status.Router())     // product-type router
		r.Mount("/tag", tag.Router())           // product-tag router
		r.Mount("/category", category.Router()) // product-category router
	})

	return r
}
