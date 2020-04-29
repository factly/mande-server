package membership

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// detail - Get membership by id
// @Summary Show a membership by id
// @Description Get membership by ID
// @Tags Membership
// @ID get-membership-by-id
// @Produce  json
// @Param id path string true "Membership ID"
// @Success 200 {object} model.Membership
// @Failure 400 {array} string
// @Router /memberships/{id} [get]
func detail(w http.ResponseWriter, r *http.Request) {

	membershipID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(membershipID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Membership{
		ID: uint(id),
	}

	err = model.DB.Model(&model.Membership{}).First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Model(&req).Association("User").Find(&req.User)
	model.DB.Model(&req).Association("Plan").Find(&req.Plan)
	model.DB.Model(&req).Association("Payment").Find(&req.Payment)
	model.DB.Model(&req.Payment).Association("Currency").Find(&req.Payment.Currency)
	json.NewEncoder(w).Encode(req)
}
