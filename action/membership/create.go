package membership

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
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
// @Success 201 {object} model.Membership
// @Failure 400 {array} string
// @Router /memberships [post]
func create(w http.ResponseWriter, r *http.Request) {

	membership := &membership{}
	json.NewDecoder(r.Body).Decode(&membership)

	validate := validator.New()
	err := validate.StructExcept(membership, "User", "Plan", "Payment")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	result := &model.Membership{
		Status:    membership.Status,
		UserID:    membership.UserID,
		PaymentID: membership.PaymentID,
		PlanID:    membership.PlanID,
	}

	err = model.DB.Model(&model.Membership{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Preload("User").Preload("Plan").Preload("Payment").Preload("Payment.Currency").First(&result)

	render.JSON(w, http.StatusCreated, result)
}
