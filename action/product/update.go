package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// update - Update product by id
// @Summary Update a product by id
// @Description Update product by ID
// @Tags Product
// @ID update-product-by-id
// @Produce json
// @Consume json
// @Param product_id path string true "Product ID"
// @Param Product body product false "Product"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products/{product_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	product := &product{}
	json.NewDecoder(r.Body).Decode(&product)

	result := &model.Product{}
	result.ID = uint(id)
	result.Tags = make([]model.Tag, 0)
	result.Datasets = make([]model.Dataset, 0)

	// check record exist or not
	err = model.DB.Model(&model.Product{}).Preload("Currency").Preload("FeaturedMedium").Preload("Tags").Preload("Datasets").First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	oldTags := result.Tags
	newTags := make([]model.Tag, 0)
	model.DB.Model(&model.Tag{}).Where(product.TagIDs).Find(&newTags)

	oldDatasets := result.Datasets
	newDatasets := make([]model.Dataset, 0)
	model.DB.Model(&model.Dataset{}).Where(product.DatasetIDs).Find(&newDatasets)

	if len(oldTags) > 0 {
		model.DB.Model(&result).Association("Tags").Delete(oldTags)
	}
	if len(oldDatasets) > 0 {
		model.DB.Model(&result).Association("Datasets").Delete(oldDatasets)
	}

	if len(newTags) == 0 {
		newTags = nil
	}
	if len(newDatasets) == 0 {
		newDatasets = nil
	}

	model.DB.Model(&result).Set("gorm:association_autoupdate", false).Updates(&model.Product{
		CurrencyID: product.CurrencyID,
		Status:     product.Status,
		Title:      product.Title,
		Price:      product.Price,
		Slug:       product.Slug,
		Tags:       newTags,
		Datasets:   newDatasets,
	}).Preload("Currency").Preload("FeaturedMedium").Preload("Tags").Preload("Datasets").First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
