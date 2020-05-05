package item

import "github.com/go-chi/chi"

// OrderItem request body
type orderItem struct {
	ExtraInfo string `json:"extra_info"`
	ProductID uint   `json:"product_id"`
}

// Router - Group of currency router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /orders/{order_id}/order-items - return list of order items
	r.Post("/", create) // POST /order-items - create a new order item and persist it

	r.Route("/{item_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /orders/{order_id}/order-items/{item_id} - read a single order item by :id
		r.Put("/", update)    // PUT /orders/{order_id}/order-items/{item_id} - update a single order item by :id
		r.Delete("/", delete) // DELETE /orders/{order_id}/order-items/{item_id} - delete a single order item by :id
	})

	return r
}
