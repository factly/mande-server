package product

import (
	"github.com/go-chi/chi"
)

type product struct {
	Title            string `json:"title" validate:"required"`
	Slug             string `json:"slug" validate:"required"`
	Price            int    `json:"price" validate:"required"`
	Status           string `json:"status"`
	CurrencyID       uint   `json:"currency_id"`
	FeaturedMediumID uint   `json:"featured_medium_id"`
	DatasetIDs       []uint `json:"dataset_ids"`
	TagIDs           []uint `json:"tag_ids"`
}

// Router - Group of product router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/", create) // POST /products - create a new product
	r.Get("/", list)    // GET /products - return list of products
	r.Route("/{product_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /products/{product_id} - read a single product by :payment_id
		r.Delete("/", delete) // DELETE /products/{product_id} - delete a single product by :product_id
		r.Put("/", update)    // PUT /products/{product_id} - update a single product by :product_id
	})

	return r
}
