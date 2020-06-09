package membership

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// details - Get membership by id
// @Summary Show a membership by id
// @Description Get membership by ID
// @Tags Membership
// @ID get-membership-by-id
// @Produce  json
// @Param membership_id path string true "Membership ID"
// @Success 200 {object} model.Membership
// @Failure 400 {array} string
// @Router /memberships/{membership_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	membershipID := chi.URLParam(r, "membership_id")
	id, err := strconv.Atoi(membershipID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	result := &model.Membership{}
	result.ID = uint(id)

	err = model.DB.Model(&model.Membership{}).First(&result).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Preload("User").Preload("Plan").Preload("Payment").Preload("Payment.Currency").First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
