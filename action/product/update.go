package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/array"
	"github.com/factly/x/errorx"
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
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	product := &product{}
	productDatasets := []model.ProductDataset{}
	productTags := []model.ProductTag{}
	json.NewDecoder(r.Body).Decode(&product)

	result := &productData{}
	result.ID = uint(id)
	result.Tags = make([]model.Tag, 0)
	result.Datasets = make([]model.Dataset, 0)

	model.DB.Model(&result.Product).Updates(&model.Product{
		CurrencyID: product.CurrencyID,
		Status:     product.Status,
		Title:      product.Title,
		Price:      product.Price,
		Slug:       product.Slug,
	}).Preload("Currency").Preload("FeaturedMedium").First(&result.Product)

	// fetch all datasets
	model.DB.Model(&model.ProductDataset{}).Where(&model.ProductDataset{
		ProductID: uint(id),
	}).Preload("Dataset").Find(&productDatasets)

	// fetch all tags
	model.DB.Model(&model.ProductTag{}).Where(&model.ProductTag{
		ProductID: uint(id),
	}).Preload("Tag").Find(&productTags)

	prevTagIDs := make([]uint, 0)
	productTagIDs := make([]uint, 0)
	mapperProductTag := map[uint]model.ProductTag{}

	for _, productTag := range productTags {
		mapperProductTag[productTag.TagID] = productTag
		prevTagIDs = append(prevTagIDs, productTag.TagID)
	}

	toCreateIDs, toDeleteIDs := array.Difference(prevTagIDs, product.TagIDs)

	// map product tag ids
	for _, id := range toDeleteIDs {
		productTagIDs = append(productTagIDs, mapperProductTag[id].ID)
	}

	// delete product tags
	if len(productTagIDs) > 0 {
		model.DB.Where(productTagIDs).Delete(model.ProductTag{})
	}

	// create product tags
	for _, id := range toCreateIDs {
		productTag := &model.ProductTag{}
		productTag.TagID = uint(id)
		productTag.ProductID = result.ID

		err = model.DB.Model(&model.ProductTag{}).Create(&productTag).Error
		if err != nil {
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	// fetch updated product tags
	updatedProductTags := []model.ProductTag{}
	model.DB.Model(&model.ProductTag{}).Where(&model.ProductTag{
		ProductID: uint(id),
	}).Preload("Tag").Find(&updatedProductTags)

	// appending previous product tags to result
	for _, productTag := range updatedProductTags {
		result.Tags = append(result.Tags, productTag.Tag)
	}

	prevDatasetIDs := make([]uint, 0)
	productDatasetIDs := make([]uint, 0)
	mapperProductDataset := map[uint]model.ProductDataset{}

	for _, productDataset := range productDatasets {
		mapperProductDataset[productDataset.DatasetID] = productDataset
		prevDatasetIDs = append(prevDatasetIDs, productDataset.DatasetID)
	}

	toCreateIDs, toDeleteIDs = array.Difference(prevDatasetIDs, product.DatasetIDs)

	// map product datasets ids
	for _, id := range toDeleteIDs {
		productDatasetIDs = append(productDatasetIDs, mapperProductDataset[id].ID)
	}

	// delete product datasets
	if len(productDatasetIDs) > 0 {
		model.DB.Where(productDatasetIDs).Delete(model.ProductDataset{})
	}

	// creating new datasets
	for _, id := range product.DatasetIDs {
		productDataset := &model.ProductDataset{}
		productDataset.DatasetID = uint(id)
		productDataset.ProductID = result.ID

		err = model.DB.Model(&model.ProductDataset{}).Create(&productDataset).Error
		if err != nil {
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	// fetch updated product datasets
	updatedProductDatasets := []model.ProductDataset{}
	model.DB.Model(&model.ProductDataset{}).Preload("Dataset").First(&updatedProductDatasets)

	// appending previous product datasets to result
	for _, productDataset := range updatedProductDatasets {
		result.Datasets = append(result.Datasets, productDataset.Dataset)
	}

	renderx.JSON(w, http.StatusOK, result)
}
