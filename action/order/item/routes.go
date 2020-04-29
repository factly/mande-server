package item

import "github.com/go-chi/chi"

// OrderItem request body
type orderItem struct {
	ExtraInfo string `json:"extra_info"`
	ProductID uint   `json:"product_id"`
	OrderID   uint   `json:"order_id"`
}

// Router - Group of currency router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", getOrderItems)    // GET /order-items - return list of order items
	r.Post("/", createOrderItem) // POST /order-items - create a new order item and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", getOrderItemByID)   // GET /order-items/{id} - read a single order item by :id
		r.Put("/", updateOrderItem)    // PUT /order-items/{id} - update a single order item by :id
		r.Delete("/", deleteOrderItem) // DELETE /order-items/{id} - delete a single order item by :id
	})

	return r
}
