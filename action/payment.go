package action

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// payment request body
type payment struct {
	Amount     int    `json:"amount"`
	Gateway    string `json:"gateway"`
	CurrencyID uint   `json:"currency_id"`
	Status     string `json:"status"`
}

// GetPayments - Get all payments
// @Summary Show all payments
// @Description Get all payments
// @Tags Payment
// @ID get-all-payments
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Payment
// @Router /payments [get]
func GetPayments(w http.ResponseWriter, r *http.Request) {

	var payments []model.Payment
	p := r.URL.Query().Get("page")
	pg, _ := strconv.Atoi(p) // pg contains page number
	l := r.URL.Query().Get("limit")
	li, _ := strconv.Atoi(l) // li contains perPage number

	offset := 0 // no. of records to skip
	limit := 5  // limt

	if li > 0 && li <= 10 {
		limit = li
	}

	if pg > 1 {
		offset = (pg - 1) * limit
	}

	model.DB.Offset(offset).Limit(limit).Preload("Currency").Model(&model.Payment{}).Find(&payments)

	json.NewEncoder(w).Encode(payments)
}

// GetPayment - Get payment by id
// @Summary Show a payment by id
// @Description Get payment by ID
// @Tags Payment
// @ID get-payment-by-id
// @Produce  json
// @Param id path string true "Payment ID"
// @Success 200 {object} model.Payment
// @Failure 400 {array} string
// @Router /payments/{id} [get]
func GetPayment(w http.ResponseWriter, r *http.Request) {

	paymentID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paymentID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Payment{
		ID: uint(id),
	}

	err = model.DB.Model(&model.Payment{}).First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Model(&req).Association("Currency").Find(&req.Currency)
	json.NewEncoder(w).Encode(req)
}

// CreatePayment - Create payment
// @Summary Create payment
// @Description Create payment
// @Tags Payment
// @ID add-payment
// @Consume json
// @Produce  json
// @Param Payment body payment true "Payment object"
// @Success 200 {object} model.Payment
// @Failure 400 {array} string
// @Router /payments [post]
func CreatePayment(w http.ResponseWriter, r *http.Request) {

	req := &model.Payment{}
	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.StructExcept(req, "Currency")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Payment{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&req).Association("Currency").Find(&req.Currency)
	json.NewEncoder(w).Encode(req)
}

// UpdatePayment - Update payment by id
// @Summary Update a payment by id
// @Description Update payment by ID
// @Tags Payment
// @ID update-payment-by-id
// @Produce json
// @Consume json
// @Param id path string true "Payment ID"
// @Param Payment body payment false "Payment"
// @Success 200 {object} model.Payment
// @Failure 400 {array} string
// @Router /payments/{id} [put]
func UpdatePayment(w http.ResponseWriter, r *http.Request) {

	paymentID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paymentID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	payment := &model.Payment{
		ID: uint(id),
	}

	req := &model.Payment{}

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&payment).Updates(&model.Payment{
		Amount:     req.Amount,
		Gateway:    req.Gateway,
		Status:     req.Status,
		CurrencyID: req.CurrencyID,
	})
	model.DB.First(&payment)
	model.DB.Model(&payment).Association("Currency").Find(&payment.Currency)
	json.NewEncoder(w).Encode(payment)
}

// DeletePayment - Delete payment by id
// @Summary Delete a payment
// @Description Delete payment by ID
// @Tags Payment
// @ID delete-payment-by-id
// @Consume  json
// @Param id path string true "Payment ID"
// @Success 200 {object} model.Payment
// @Failure 400 {array} string
// @Router /payments/{id} [delete]
func DeletePayment(w http.ResponseWriter, r *http.Request) {

	paymentID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paymentID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	payment := &model.Payment{
		ID: uint(id),
	}

	// check record exists or not
	err = model.DB.First(&payment).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Model(&payment).Association("Currency").Find(&payment.Currency)
	model.DB.Delete(&payment)

	json.NewEncoder(w).Encode(payment)
}
