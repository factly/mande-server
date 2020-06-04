package payment

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
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

	err := validation.Validator.Struct(payment)
	if err != nil {
		validation.ValidatorErrors(w, r, err)
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

	render.JSON(w, http.StatusCreated, result)
}
