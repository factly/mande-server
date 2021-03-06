package plan

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/meilisearchx"
	"github.com/factly/x/middlewarex"
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
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param Plan body plan true "Plan object"
// @Success 201 {object} model.Plan
// @Router /plans [post]
func Create(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
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

	result := &model.Plan{
		Name:        plan.Name,
		Description: plan.Description,
		Duration:    plan.Duration,
		Status:      plan.Status,
		CurrencyID:  plan.CurrencyID,
		Price:       plan.Price,
		AllProducts: plan.AllProducts,
		Users:       plan.Users,
	}

	result.Catalogs = make([]model.Catalog, 0)

	model.DB.Model(&model.Catalog{}).Where(plan.CatalogIDs).Find(&result.Catalogs)

	tx := model.DB.WithContext(context.WithValue(r.Context(), userContext, uID)).Begin()
	err = tx.Model(&model.Plan{}).Create(&result).Error

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	if !plan.AllProducts {
		tx.Preload("Currency").Preload("Catalogs").Preload("Catalogs.Products").Preload("Catalogs.Products.Currency").Preload("Catalogs.Products.Datasets").Preload("Catalogs.Products.Tags").First(&result)
	} else {
		tx.Preload("Currency").First(&result)
	}

	// Insert into meili index
	meiliObj := map[string]interface{}{
		"id":           result.ID,
		"kind":         "plan",
		"name":         result.Name,
		"description":  result.Description,
		"duration":     result.Duration,
		"status":       result.Status,
		"catalog_ids":  plan.CatalogIDs,
		"all_products": plan.AllProducts,
		"users":        plan.Users,
	}

	err = meilisearchx.AddDocument("data-portal", meiliObj)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()
	renderx.JSON(w, http.StatusCreated, result)
}
