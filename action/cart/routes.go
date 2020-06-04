package cart

import (
	"github.com/factly/data-portal-server/action/cart/item"
	"github.com/go-chi/chi"
)

// Cart request body
type cart struct {
	Status string `json:"status" validate:"required"`
	UserID uint   `json:"user_id" validate:"required"`
}

// Router - Group of cart router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /carts - return list of carts
	r.Post("/", create) // POST /carts - create a new cart and persist it

	r.Route("/{cart_id}", func(r chi.Router) {
		r.Get("/", details)              // GET /carts/{cart_id} - read a single cart by :cart_id
		r.Put("/", update)               // PUT /carts/{cart_id} - update a single cart by :cart_id
		r.Delete("/", delete)            // DELETE /carts/{cart_id} - delete a single cart by :cart_id
		r.Mount("/items", item.Router()) // cart-item router
	})

	return r
}
