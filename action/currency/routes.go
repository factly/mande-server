package currency

import (
	"github.com/factly/mande-server/model"
	"github.com/go-chi/chi"
)

// currency request body
type currency struct {
	IsoCode string `json:"iso_code" validate:"required"`
	Name    string `json:"name" validate:"required"`
}

var userContext model.ContextKey = "currency_user"

// UserRouter - Group of currency router
func UserRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list) // GET /currencies - return list of currencies

	r.Route("/{currency_id}", func(r chi.Router) {
		r.Get("/", details) // GET /currencies/{currency_id} - read a single currency by :currency_id
	})

	return r
}

// AdminRouter - Group of currency router
func AdminRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)                  // GET /currencies - return list of currencies
	r.Post("/", create)               // POST /currencies - create a new currency and persist it
	r.Post("/default", createDefault) // POST /currencies/default - create a new currency and persist it

	r.Route("/{currency_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /currencies/{currency_id} - read a single currency by :currency_id
		r.Put("/", update)    // PUT /currencies/{currency_id} - update a single currency by :currency_id
		r.Delete("/", delete) // DELETE /currencies/{currency_id} - delete a single currency by :currency_id
	})

	return r
}
