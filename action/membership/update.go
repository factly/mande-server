package membership

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// update - Update membership by id
// @Summary Update a membership by id
// @Description Update membership by ID
// @Tags Membership
// @ID update-membership-by-id
// @Produce json
// @Consume json
// @Param id path string true "Membership ID"
// @Param Membership body membership false "Membership"
// @Success 200 {object} model.Membership
// @Failure 400 {array} string
// @Router /memberships/{id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	membershipID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(membershipID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	membership := &membership{}
	result := &model.Membership{}
	result.ID = uint(id)

	json.NewDecoder(r.Body).Decode(&membership)

	model.DB.Model(&result).Updates(model.Membership{
		UserID:    membership.UserID,
		PaymentID: membership.PaymentID,
		PlanID:    membership.PlanID,
		Status:    membership.Status,
	}).First(&result)

	model.DB.Model(&result).Association("User").Find(&result.User)
	model.DB.Model(&result).Association("Plan").Find(&result.Plan)
	model.DB.Model(&result).Association("Payment").Find(&result.Payment)
	model.DB.Model(&result.Payment).Association("Currency").Find(&result.Payment.Currency)

	render.JSON(w, http.StatusOK, result)
}
