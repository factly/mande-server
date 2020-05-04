package membership

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// details - Get membership by id
// @Summary Show a membership by id
// @Description Get membership by ID
// @Tags Membership
// @ID get-membership-by-id
// @Produce  json
// @Param id path string true "Membership ID"
// @Success 200 {object} model.Membership
// @Failure 400 {array} string
// @Router /memberships/{id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	membershipID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(membershipID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	membership := &model.Membership{}
	membership.ID = uint(id)

	err = model.DB.Model(&model.Membership{}).First(&membership).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Model(&membership).Association("User").Find(&membership.User)
	model.DB.Model(&membership).Association("Plan").Find(&membership.Plan)
	model.DB.Model(&membership).Association("Payment").Find(&membership.Payment)
	model.DB.Model(&membership.Payment).Association("Currency").Find(&membership.Payment.Currency)

	util.Render(w, http.StatusOK, membership)
}
