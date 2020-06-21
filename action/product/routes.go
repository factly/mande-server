package product

import (
	"github.com/factly/data-portal-server/model"
	"github.com/go-chi/chi"
)

type product struct {
	Title           string `json:"title" validate:"required"`
	Slug            string `json:"slug" validate:"required"`
	Price           int    `json:"price" validate:"required"`
	Status          string `json:"status"`
	CurrencyID      uint   `json:"currency_id"`
	FeaturedMediaID uint   `json:"featured_media_id"`
	DatasetIDs      []uint `json:"dataset_ids"`
	TagIDs          []uint `json:"tag_ids"`
}

type productData struct {
	model.Product
	Tags     []model.Tag     `json:"tags"`
	Datasets []model.Dataset `json:"datasets"`
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
