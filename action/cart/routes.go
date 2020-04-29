package cart

import "github.com/go-chi/chi"

// Cart request body
type cart struct {
	Status string `json:"status"`
	UserID uint   `json:"user_id"`
}

// Router - Group of currency router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", getCarts)    // GET /carts - return list of carts
	r.Post("/", createCart) // POST /carts - create a new cart and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", getCartByID)   // GET /carts/{id} - read a single cart by :id
		r.Put("/", updateCart)    // PUT /carts/{id} - update a single cart by :id
		r.Delete("/", deleteCart) // DELETE /carts/{id} - delete a single cart by :id
	})

	return r
}
