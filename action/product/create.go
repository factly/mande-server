package product

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
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
// @Param Product body product true "Product object"
// @Success 201 {object} model.Product
// @Failure 400 {array} string
// @Router /products [post]
func create(w http.ResponseWriter, r *http.Request) {

	product := &product{}
	json.NewDecoder(r.Body).Decode(&product)

	validationError := validationx.Check(product)
	if validationError != nil {
		errorx.Render(w, validationError)
		return
	}

	result := &productData{}
	result.Tags = make([]model.Tag, 0)
	result.Datasets = make([]model.Dataset, 0)
	result.Product = model.Product{
		Title:            product.Title,
		Slug:             product.Slug,
		Price:            product.Price,
		Status:           product.Status,
		CurrencyID:       product.CurrencyID,
		FeaturedMediumID: product.FeaturedMediumID,
	}

	err := model.DB.Model(&model.Product{}).Create(&result.Product).Error

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	model.DB.Preload("Currency").Preload("FeaturedMedium").First(&result.Product)

	for _, id := range product.DatasetIDs {
		productDataset := &model.ProductDataset{}
		productDataset.DatasetID = uint(id)
		productDataset.ProductID = result.ID

		err = model.DB.Model(&model.ProductDataset{}).Create(&productDataset).Error
		if err != nil {
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
		model.DB.Model(&model.ProductDataset{}).Preload("Dataset").First(&productDataset)
	}

	// fetch all product datasets
	productDatasets := []model.ProductDataset{}
	model.DB.Model(&model.ProductDataset{}).Where(&model.ProductDataset{
		ProductID: result.Product.ID,
	}).Preload("Dataset").Find(&productDatasets)

	// appending product datasets to result
	for _, productDataset := range productDatasets {
		result.Datasets = append(result.Datasets, productDataset.Dataset)
	}

	for _, id := range product.TagIDs {
		productTag := &model.ProductTag{}
		productTag.TagID = uint(id)
		productTag.ProductID = result.ID

		err = model.DB.Model(&model.ProductTag{}).Create(&productTag).Error
		if err != nil {
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	// fetch all product tags
	productTags := []model.ProductTag{}
	model.DB.Model(&model.ProductTag{}).Where(&model.ProductTag{
		ProductID: result.Product.ID,
	}).Preload("Tag").Find(&productTags)

	// appending product tags to result
	for _, productTag := range productTags {
		result.Tags = append(result.Tags, productTag.Tag)
	}

	renderx.JSON(w, http.StatusCreated, result)
}
