package product

import (
	"github.com/factly/data-portal-server/action/product/category"
	"github.com/factly/data-portal-server/action/product/prodtype"
	"github.com/factly/data-portal-server/action/product/status"
	"github.com/factly/data-portal-server/action/product/tag"
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

	r.Post("/", create) // POST /products - create a new product
	r.Get("/", list)    // GET /products - return list of products
	r.Route("/{product_id}", func(r chi.Router) {
		r.Get("/", details)                     // GET /products/{product_id} - read a single product by :payment_id
		r.Delete("/", delete)                   // DELETE /products/{product_id} - delete a single product by :product_id
		r.Put("/", update)                      // PUT /products/{product_id} - update a single product by :product_id
		r.Mount("/type", prodtype.Router())     // product-type router
		r.Mount("/status", status.Router())     // product-status router
		r.Mount("/tag", tag.Router())           // product-tag router
		r.Mount("/category", category.Router()) // product-category router
	})

	return r
}
