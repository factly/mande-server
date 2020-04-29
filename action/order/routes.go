package order

import "github.com/go-chi/chi"

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

	r.Get("/", getOrders)    // GET /orders - return list of orders
	r.Post("/", createOrder) // POST /orders - create a new order and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", getOrderByID)   // GET /orders/{id} - read a single order by :id
		r.Put("/", updateOrder)    // PUT /orders/{id} - update a single order by :id
		r.Delete("/", deleteOrder) // DELETE /orders/{id} - delete a single order by :id
	})

	return r
}
