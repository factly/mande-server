package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/factly/data-portal-api/validationerrors"
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

// GetPayment - Get payment by id
// @Summary Show a payment by id
// @Description Get payment by ID
// @Tags Payment
// @ID get-payment-by-id
// @Produce  json
// @Param id path string true "Payment ID"
// @Success 200 {object} models.Payment
// @Failure 400 {array} string
// @Router /payments/{id} [get]
func GetPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	paymentID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paymentID)

	if err != nil {
		validationerrors.InvalidID(w, r)
		return
	}

	req := &models.Payment{
		ID: uint(id),
	}

	models.DB.Model(&models.Payment{}).First(&req)

	models.DB.Model(&req).Association("Currency").Find(&req.Currency)
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
// @Success 200 {object} models.Payment
// @Failure 400 {array} string
// @Router /payments [post]
func CreatePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := &models.Payment{}
	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.StructExcept(req, "Currency")
	if err != nil {
		msg := err.Error()
		validationerrors.ValidErrors(w, r, msg)
		return
	}

	err = models.DB.Model(&models.Payment{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	models.DB.Model(&req).Association("Currency").Find(&req.Currency)
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
// @Success 200 {object} models.Payment
// @Failure 400 {array} string
// @Router /payments/{id} [put]
func UpdatePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	paymentID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paymentID)

	if err != nil {
		validationerrors.InvalidID(w, r)
		return
	}

	payment := &models.Payment{
		ID: uint(id),
	}

	req := &models.Payment{}

	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&payment).Updates(&models.Payment{
		Amount:     req.Amount,
		Gateway:    req.Gateway,
		Status:     req.Status,
		CurrencyID: req.CurrencyID,
	})
	models.DB.First(&payment)
	models.DB.Model(&payment).Association("Currency").Find(&payment.Currency)
	json.NewEncoder(w).Encode(payment)
}

// DeletePayment - Delete payment by id
// @Summary Delete a payment
// @Description Delete payment by ID
// @Tags Payment
// @ID delete-payment-by-id
// @Consume  json
// @Param id path string true "Payment ID"
// @Success 200 {object} models.Payment
// @Failure 400 {array} string
// @Router /payments/{id} [delete]
func DeletePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	paymentID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paymentID)

	if err != nil {
		validationerrors.InvalidID(w, r)
		return
	}

	payment := &models.Payment{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&payment).Error
	if err != nil {
		validationerrors.RecordNotFound(w, r)
		return
	}
	models.DB.Model(&payment).Association("Currency").Find(&payment.Currency)
	models.DB.Delete(&payment)

	json.NewEncoder(w).Encode(payment)
}
