package item

import "github.com/go-chi/chi"

// CartItem request body
type cartItem struct {
	IsDeleted bool `json:"is_deleted"`
	CartID    uint `json:"cart_id"`
	ProductID uint `json:"product_id"`
}

// Router - Group of cart-item router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /carts/{cart_id}/carts/{cart_id}/cart-items - return list of cart items
	r.Post("/", create) // POST /carts/{cart_id}/cart-items - create a new cart item and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", detail)    // GET /carts/{cart_id}/cart-items/{item_id} - read a single cart item by :id
		r.Put("/", update)    // PUT /carts/{cart_id}/cart-items/{item_id} - update a single cart item by :id
		r.Delete("/", delete) // DELETE /carts/{cart_id}/cart-items/{item_id} - delete a single cart item by :id
	})

	return r
}