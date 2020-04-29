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

	r.Post("/", createProduct)
	r.Get("/", getProducts)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", getProductByID)
		r.Delete("/", deleteProduct)
		r.Put("/", updateProduct)
		r.Post("/type", prodtype.CreateProductType)
		r.Get("/type", prodtype.GetProductTypes)
		r.Post("/status", status.CreateStatus)
		r.Get("/status", status.GetStatuses)
		r.Route("/tag", func(r chi.Router) {
			r.Post("/", tag.CreateProductTag)
			r.Route("/{tid}", func(r chi.Router) {
				r.Delete("/", tag.DeleteProductTag)
			})
		})
		r.Route("/category", func(r chi.Router) {
			r.Post("/", category.CreateProductCategory)
			r.Route("/{cid}", func(r chi.Router) {
				r.Delete("/", category.DeleteProductCategory)
			})

		})
	})

	return r
}
