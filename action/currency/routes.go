package currency

import "github.com/go-chi/chi"

// currency request body
type currency struct {
	IsoCode string `json:"iso_code"`
	Name    string `json:"name"`
}

// Router - Group of currency router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /currencies - return list of currencies
	r.Post("/", create) // POST /currencies - create a new currency and persist it

	r.Route("/{currency_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /currencies/{currency_id} - read a single currency by :currency_id
		r.Put("/", update)    // PUT /currencies/{currency_id} - update a single currency by :currency_id
		r.Delete("/", delete) // DELETE /currencies/{currency_id} - delete a single currency by :currency_id
	})

	return r
}
