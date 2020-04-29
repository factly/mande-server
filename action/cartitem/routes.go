package cartitem

import "github.com/go-chi/chi"

// CartItem request body
type cartItem struct {
	IsDeleted bool `json:"is_deleted"`
	CartID    uint `json:"cart_id"`
	ProductID uint `json:"product_id"`
}

// Router - Group of currency router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", getCartItems)    // GET /cart-items - return list of cart items
	r.Post("/", createCartItem) // POST /cart-items - create a new cart item and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", getCartItemByID)   // GET /cart-items/{id} - read a single cart item by :id
		r.Put("/", updateCartItem)    // PUT /cart-items/{id} - update a single cart item by :id
		r.Delete("/", deleteCartItem) // DELETE /cart-items/{id} - delete a single cart item by :id
	})

	return r
}
