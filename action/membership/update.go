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

	req := &membership{}
	membership := &model.Membership{}
	membership.ID = uint(id)

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&membership).Updates(model.Membership{
		UserID:    req.UserID,
		PaymentID: req.PaymentID,
		PlanID:    req.PlanID,
		Status:    req.Status,
	})

	model.DB.First(&membership)
	model.DB.Model(&membership).Association("User").Find(&membership.User)
	model.DB.Model(&membership).Association("Plan").Find(&membership.Plan)
	model.DB.Model(&membership).Association("Payment").Find(&membership.Payment)
	model.DB.Model(&membership.Payment).Association("Currency").Find(&membership.Payment.Currency)

	render.JSON(w, http.StatusOK, membership)
}
