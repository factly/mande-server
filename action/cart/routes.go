package cart

import (
	"github.com/go-chi/chi"
)

// CartItem request body
type cartitem struct {
	Status    string `json:"status" validate:"required"`
	ProductID uint   `json:"product_id" validate:"required"`
}

// Router - Group of cart router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /carts - return list of cart items
	r.Post("/", create) // POST /carts - add a new cart item

	r.Route("/{cartitem_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /carts/{cartitem_id} - get a single cart item
		r.Delete("/", delete) // DELETE /carts/{cartitem_id} - delete a cart item entry
	})

	return r
}
