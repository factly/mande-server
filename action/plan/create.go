package plan

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/meili"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// Create - create plan
// @Summary Create plan
// @Description create plan
// @Tags Plan
// @ID add-plan
// @Consume json
// @Produce  json
// @Param Plan body plan true "Plan object"
// @Success 201 {object} model.Plan
// @Router /plans [post]
func Create(w http.ResponseWriter, r *http.Request) {

	plan := &plan{}

	err := json.NewDecoder(r.Body).Decode(&plan)
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

	result := &model.Plan{
		Name:        plan.Name,
		Description: plan.Description,
		Duration:    plan.Duration,
		Status:      plan.Status,
		CurrencyID:  plan.CurrencyID,
		Price:       plan.Price,
	}

	result.Catalogs = make([]model.Catalog, 0)

	model.DB.Model(&model.Catalog{}).Where(plan.CatalogIDs).Find(&result.Catalogs)

	tx := model.DB.Begin()
	err = tx.Model(&model.Plan{}).Set("gorm:association_autoupdate", false).Create(&result).Error

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	tx.Preload("Currency").Preload("Catalogs").Preload("Catalogs.Products").Preload("Catalogs.Products.Currency").Preload("Catalogs.Products.Datasets").Preload("Catalogs.Products.Tags").First(&result)

	// Insert into meili index
	meiliObj := map[string]interface{}{
		"id":          result.ID,
		"kind":        "plan",
		"name":        result.Name,
		"description": result.Description,
		"duration":    result.Duration,
		"status":      result.Status,
		"catalog_ids": plan.CatalogIDs,
	}

	err = meili.AddDocument(meiliObj)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()
	renderx.JSON(w, http.StatusCreated, result)
}
