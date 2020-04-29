package membership

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-playground/validator/v10"
)

// create - Create membership
// @Summary Create membership
// @Description Create membership
// @Tags Membership
// @ID add-membership
// @Consume json
// @Produce  json
// @Param Membership body membership true "Membership object"
// @Success 200 {object} model.Membership
// @Failure 400 {array} string
// @Router /memberships [post]
func create(w http.ResponseWriter, r *http.Request) {

	req := &model.Membership{}
	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.StructExcept(req, "User", "Plan", "Payment")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Membership{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&req).Association("User").Find(&req.User)
	model.DB.Model(&req).Association("Plan").Find(&req.Plan)
	model.DB.Model(&req).Association("Payment").Find(&req.Payment)
	model.DB.Model(&req.Payment).Association("Currency").Find(&req.Payment.Currency)

	json.NewEncoder(w).Encode(req)
}
