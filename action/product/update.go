package product

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/mande-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/meilisearchx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// update - Update product by id
// @Summary Update a product by id
// @Description Update product by ID
// @Tags Product
// @ID update-product-by-id
// @Produce json
// @Consume json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param product_id path string true "Product ID"
// @Param Product body product false "Product"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products/{product_id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
		return
	}

	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
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

	result := model.Product{}
	result.ID = uint(id)
	result.Tags = make([]model.Tag, 0)
	result.Datasets = make([]model.Dataset, 0)

	// check record exist or not
	err = model.DB.First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	tx := model.DB.Begin()

	newTags := make([]model.Tag, 0)
	if len(product.TagIDs) > 0 {
		model.DB.Model(&model.Tag{}).Where(product.TagIDs).Find(&newTags)
		err = tx.Model(&result).Association("Tags").Replace(&newTags)
	} else {
		err = tx.Model(&result).Association("Tags").Clear()
	}

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	newDatasets := make([]model.Dataset, 0)
	if len(product.DatasetIDs) > 0 {
		model.DB.Model(&model.Dataset{}).Where(product.DatasetIDs).Find(&newDatasets)
		err = tx.Model(&result).Association("Datasets").Replace(&newDatasets)
	} else {
		err = tx.Model(&result).Association("Datasets").Clear()
	}

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	featuredMediumID := &product.FeaturedMediumID
	if product.FeaturedMediumID == 0 {
		err = tx.Omit("Datasets", "Tags").Model(result).Updates(map[string]interface{}{"featured_medium_id": nil}).Error
		featuredMediumID = nil
		if err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	err = tx.Model(&result).Omit("Datasets", "Tags").Updates(&model.Product{
		Base:             model.Base{UpdatedByID: uint(uID)},
		CurrencyID:       product.CurrencyID,
		Status:           product.Status,
		Title:            product.Title,
		Price:            product.Price,
		FeaturedMediumID: featuredMediumID,
		Slug:             product.Slug,
		Description:      product.Description,
	}).Preload("Currency").Preload("FeaturedMedium").Preload("Tags").Preload("Datasets").First(&result).Error

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	// Update into meili index
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
		"description": result.Description,
	}

	err = meilisearchx.UpdateDocument("data-portal", meiliObj)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()

	renderx.JSON(w, http.StatusOK, result)
}
