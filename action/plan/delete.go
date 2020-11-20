package plan

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/meili"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete plan by id
// @Summary Delete a plan
// @Description Delete plan by ID
// @Tags Plan
// @ID delete-plan-by-id
// @Consume  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param plan_id path string true "Plan ID"
// @Success 200
// @Failure 400 {array} string
// @Router /plans/{plan_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	planID := chi.URLParam(r, "plan_id")
	id, err := strconv.Atoi(planID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Plan{}

	result.ID = uint(id)

	// check record exists or not
	err = model.DB.Preload("Catalogs").First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	// check if plan is associated with membership
	var totAssociated int64
	model.DB.Model(&model.Membership{}).Where(&model.Membership{
		PlanID: uint(id),
	}).Count(&totAssociated)

	if totAssociated != 0 {
		loggerx.Error(errors.New("plan is associated with membership"))
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	tx := model.DB.Begin()

	// Delete all associations
	_ = tx.Model(&result).Association("Catalogs").Delete(result.Catalogs)

	tx.Delete(&result)

	err = meili.DeleteDocument(result.ID, "plan")
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()
	renderx.JSON(w, http.StatusOK, nil)
}
