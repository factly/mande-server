package payment

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-playground/validator/v10"
)

// create - Create payment
// @Summary Create payment
// @Description Create payment
// @Tags Payment
// @ID add-payment
// @Consume json
// @Produce  json
// @Param Payment body payment true "Payment object"
// @Success 201 {object} model.Payment
// @Failure 400 {array} string
// @Router /payments [post]
func create(w http.ResponseWriter, r *http.Request) {

	payment := &payment{}
	json.NewDecoder(r.Body).Decode(&payment)

	validate := validator.New()
	err := validate.StructExcept(payment, "Currency")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	result := &model.Payment{
		Amount:     payment.Amount,
		Gateway:    payment.Gateway,
		CurrencyID: payment.CurrencyID,
		Status:     payment.Status,
	}

	err = model.DB.Model(&model.Payment{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&result).Preload("Currency").Find(&result.Currency)

	renderx.JSON(w, http.StatusCreated, result)
}
