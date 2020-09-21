package order

import (
	"github.com/go-chi/chi"
)

// Router - Group of order router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /orders - return list of orders
	r.Post("/", create) // POST /orders - create a new order and persist it

	r.Route("/{order_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /orders/{order_id} - read a single order by :order_id
		r.Delete("/", delete) // DELETE /orders/{order_id} - delete a single order by :order_id
	})

	return r
}
