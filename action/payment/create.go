package payment

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/validation"
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
// @Success 200 {object} model.Payment
// @Failure 400 {array} string
// @Router /payments [post]
func create(w http.ResponseWriter, r *http.Request) {

	payment := &model.Payment{}
	json.NewDecoder(r.Body).Decode(&payment)

	validate := validator.New()
	err := validate.StructExcept(payment, "Currency")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Payment{}).Create(&payment).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&payment).Association("Currency").Find(&payment.Currency)

	util.Render(w, http.StatusOK, payment)
}
