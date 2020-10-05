package payment

import "github.com/go-chi/chi"

// payment request body
type payment struct {
	Gateway           string `json:"gateway" `
	CurrencyID        uint   `json:"currency_id" validate:"required"`
	Status            string `json:"status"`
	For               string `json:"for" validate:"required"`
	EntityID          uint   `json:"entity_id" validate:"required"`
	RazorpayPaymentID string `json:"razorpay_payment_id" validate:"required"`
	RazorpaySignature string `json:"razorpay_signature" validate:"required"`
}

// UserRouter - Group of payment router
func UserRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /payments - return list of payments
	r.Post("/", create) // POST /payments - create a new payment and persist it

	r.Route("/{payment_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /payments/{payment_id} - read a single payment by :payment_id
		r.Delete("/", delete) // DELETE /payments/{payment_id} - delete a single payment by :payment_id
	})

	return r
}

// AdminRouter - Group of payment router
func AdminRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list) // GET /payments - return list of payments

	r.Route("/{payment_id}", func(r chi.Router) {
		r.Get("/", details) // GET /payments/{payment_id} - read a single payment by :payment_id
	})

	return r
}
