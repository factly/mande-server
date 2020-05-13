package order

import (
	"github.com/factly/data-portal-server/action/order/item"
	"github.com/go-chi/chi"
)

// Order request body
type order struct {
	UserID    uint   `json:"user_id"`
	Status    string `json:"status"`
	PaymentID uint   `json:"payment_id"`
	CartID    uint   `json:"cart_id"`
}

// Router - Group of order router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /orders - return list of orders
	r.Post("/", create) // POST /orders - create a new order and persist it

	r.Route("/{order_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /orders/{order_id} - read a single order by :order_id
		r.Put("/", update)    // PUT /orders/{order_id} - update a single order by :order_id
		r.Delete("/", delete) // DELETE /orders/{order_id} - delete a single order by :order_id
		r.Mount("/items", item.Router())
	})

	return r
}
