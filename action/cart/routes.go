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

	r.Get("/", list)    // GET /carts - return list of carts
	r.Post("/", create) // POST /carts - create a new cart and persist it

	r.Route("/{cart_id}", func(r chi.Router) {
		r.Get("/", detail)    // GET /carts/{cart_id} - read a single cart by :id
		r.Put("/", update)    // PUT /carts/{cart_id} - update a single cart by :id
		r.Delete("/", delete) // DELETE /carts/{cart_id} - delete a single cart by :id
	})

	return r
}
