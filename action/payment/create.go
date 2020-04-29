package payment

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-playground/validator/v10"
)

// createPayment - Create payment
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
func createPayment(w http.ResponseWriter, r *http.Request) {

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
