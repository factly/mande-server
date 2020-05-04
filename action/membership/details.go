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

	req := &model.Membership{}
	req.ID = uint(id)

	err = model.DB.Model(&model.Membership{}).First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Model(&req).Association("User").Find(&req.User)
	model.DB.Model(&req).Association("Plan").Find(&req.Plan)
	model.DB.Model(&req).Association("Payment").Find(&req.Payment)
	model.DB.Model(&req.Payment).Association("Currency").Find(&req.Payment.Currency)

	util.Render(w, http.StatusOK, req)
}
