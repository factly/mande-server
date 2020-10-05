package plan

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/meili"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// update - Update plan by id
// @Summary Update a plan by id
// @Description Update plan by ID
// @Tags Plan
// @ID update-plan-by-id
// @Produce json
// @Consume json
// @Param plan_id path string true "Plan ID"
// @Param Plan body plan false "Plan"
// @Success 200 {object} model.Plan
// @Failure 400 {array} string
// @Router /plans/{plan_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	planID := chi.URLParam(r, "plan_id")
	id, err := strconv.Atoi(planID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	plan := &plan{}

	err = json.NewDecoder(r.Body).Decode(&plan)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(plan)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := model.Plan{}
	result.ID = uint(id)

	err = model.DB.Preload("Catalogs").First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	tx := model.DB.Begin()

	oldCatalogs := result.Catalogs
	newCatalogs := make([]model.Catalog, 0)

	if len(plan.CatalogIDs) > 0 {
		model.DB.Model(&model.Catalog{}).Where(plan.CatalogIDs).Find(&newCatalogs)
	}

	if len(oldCatalogs) > 0 {
		if err := tx.Model(&result).Association("Catalogs").Delete(oldCatalogs).Error; err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}

	}
	tx.Model(&result).Set("gorm:association_autoupdate", false).Updates(model.Plan{
		Name:        plan.Name,
		Description: plan.Description,
		Duration:    plan.Duration,
		Status:      plan.Status,
		Price:       plan.Price,
		CurrencyID:  plan.CurrencyID,
		Catalogs:    newCatalogs,
	}).Preload("Catalogs").Preload("Catalogs.Products").Preload("Catalogs.Products.Currency").Preload("Catalogs.Products.Datasets").Preload("Catalogs.Products.Tags").First(&result)

	// Update into meili index
	meiliObj := map[string]interface{}{
		"id":          result.ID,
		"kind":        "plan",
		"name":        result.Name,
		"description": result.Description,
		"duration":    result.Duration,
		"status":      result.Status,
		"catalog_ids": plan.CatalogIDs,
	}

	err = meili.UpdateDocument(meiliObj)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()
	renderx.JSON(w, http.StatusOK, result)
}
