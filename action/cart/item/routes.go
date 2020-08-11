package item

import "github.com/go-chi/chi"

// CartItem request body
type cartItem struct {
	ProductID uint `json:"product_id" validate:"required"`
}

// Router - Group of cart-item router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /carts/{cart_id}/cart-items - return list of cart items
	r.Post("/", create) // POST /carts/{cart_id}/cart-items - create a new cart item and persist it

	r.Route("/{item_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /carts/{cart_id}/cart-items/{item_id} - read a single cart item by :item_id
		r.Delete("/", delete) // DELETE /carts/{cart_id}/cart-items/{item_id} - delete a single cart item by :item_id
	})

	return r
}
