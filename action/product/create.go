package product

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/util/meili"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - Create product
// @Summary Create product
// @Description Create product
// @Tags Product
// @ID add-product
// @Consume json
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param Product body product true "Product object"
// @Success 201 {object} model.Product
// @Failure 400 {array} string
// @Router /products [post]
func create(w http.ResponseWriter, r *http.Request) {
	uID, err := util.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
		return
	}

	product := &product{}
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(product)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	featuredMediumID := &product.FeaturedMediumID
	if product.FeaturedMediumID == 0 {
		featuredMediumID = nil
	}

	result := model.Product{}
	result.Tags = make([]model.Tag, 0)
	result.Datasets = make([]model.Dataset, 0)
	result = model.Product{
		Title:            product.Title,
		Slug:             product.Slug,
		Price:            product.Price,
		Status:           product.Status,
		CurrencyID:       product.CurrencyID,
		FeaturedMediumID: featuredMediumID,
	}

	if len(product.TagIDs) > 0 {
		model.DB.Model(&model.Tag{}).Where(product.TagIDs).Find(&result.Tags)
	}

	if len(product.DatasetIDs) > 0 {
		model.DB.Model(&model.Dataset{}).Where(product.DatasetIDs).Find(&result.Datasets)
	}

	tx := model.DB.WithContext(context.WithValue(r.Context(), userContext, uID)).Begin()
	err = tx.Model(&model.Product{}).Create(&result).Error

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	tx.Preload("Currency").Preload("FeaturedMedium").Preload("Tags").Preload("Datasets").First(&result)

	// Insert into meili index
	meiliObj := map[string]interface{}{
		"id":          result.ID,
		"kind":        "product",
		"title":       result.Title,
		"slug":        result.Slug,
		"price":       result.Price,
		"status":      result.Status,
		"currency_id": result.CurrencyID,
		"tag_ids":     product.TagIDs,
		"dataset_ids": product.DatasetIDs,
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
