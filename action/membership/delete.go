package membership

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/meilisearchx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete membership by id
// @Summary Delete a membership
// @Description Delete membership by ID
// @Tags Membership
// @ID delete-membership-by-id
// @Consume  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param membership_id path string true "Membership ID"
// @Success 200
// @Failure 400 {array} string
// @Router /memberships/{membership_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	membershipID := chi.URLParam(r, "membership_id")
	id, err := strconv.Atoi(membershipID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Membership{}
	result.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}
	tx := model.DB.Begin()
	tx.Delete(&result)

	err = meilisearchx.DeleteDocument("data-portal", result.ID, "membership")
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()
	renderx.JSON(w, http.StatusOK, nil)
}
