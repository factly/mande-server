package prodtype

import "github.com/go-chi/chi"

// productType request body
type productType struct {
	Name string `json:"name"`
}

// Router - Group of product-type router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/", create) // POST /types - create a new type and persist it
	r.Get("/", list)    // GET /types - return list of types

	r.Route("/{type_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /types/{type_id} - read a single type by :type_id
		r.Put("/", update)    // PUT /types/{type_id} - update a single type by :type_id
		r.Delete("/", delete) // DELETE /types/{type_id} - delete a single type by :type_id
	})
	return r
}
