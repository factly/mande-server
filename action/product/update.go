package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
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
		validation.InvalidID(w, r)
		return
	}

	product := &product{}
	datasets := []model.ProductDataset{}
	tags := []model.ProductTag{}
	json.NewDecoder(r.Body).Decode(&product)

	result := &productData{}
	result.ID = uint(id)

	model.DB.Model(&result.Product).Updates(&model.Product{
		CurrencyID: product.CurrencyID,
		Status:     product.Status,
		Title:      product.Title,
		Price:      product.Price,
		Slug:       product.Slug,
	}).Preload("Currency").Preload("FeaturedMedia").First(&result.Product)

	// fetch all datasets
	model.DB.Model(&model.ProductDataset{}).Where(&model.ProductDataset{
		ProductID: uint(id),
	}).Preload("Dataset").Find(&datasets)

	// fetch all tags
	model.DB.Model(&model.ProductTag{}).Where(&model.ProductTag{
		ProductID: uint(id),
	}).Preload("Tag").Find(&tags)

	// delete tags
	for _, t := range tags {
		present := false
		for _, id := range product.TagIDs {
			if t.TagID == id {
				present = true
			}
		}
		if present == false {
			model.DB.Where(&model.ProductTag{
				TagID:     t.TagID,
				ProductID: uint(id),
			}).Delete(model.ProductTag{})
		}
	}

	// creating new tags
	for _, id := range product.TagIDs {
		present := false
		for _, t := range tags {
			if t.TagID == id {
				present = true
				result.Tags = append(result.Tags, t.Tag)
			}
		}
		if present == false {
			productTag := &model.ProductTag{}
			productTag.TagID = uint(id)
			productTag.ProductID = result.ID

			err = model.DB.Model(&model.ProductTag{}).Create(&productTag).Error

			if err != nil {
				return
			}
			model.DB.Model(&model.ProductTag{}).Preload("Tag").First(&productTag)
			result.Tags = append(result.Tags, productTag.Tag)
		}
	}

	// delete datasets
	for _, d := range datasets {
		present := false
		for _, id := range product.DatasetIDs {
			if d.DatasetID == id {
				present = true
			}
		}
		if present == false {
			model.DB.Where(&model.ProductDataset{
				DatasetID: d.DatasetID,
				ProductID: uint(id),
			}).Delete(model.ProductDataset{})
		}
	}

	// creating new datasets
	for _, id := range product.DatasetIDs {
		present := false
		for _, d := range datasets {
			if d.DatasetID == id {
				present = true
				result.Datasets = append(result.Datasets, d.Dataset)
			}
		}
		if present == false {
			productDataset := &model.ProductDataset{}
			productDataset.DatasetID = uint(id)
			productDataset.ProductID = result.ID

			err = model.DB.Model(&model.ProductDataset{}).Create(&productDataset).Error

			if err != nil {
				return
			}

			model.DB.Model(&model.ProductDataset{}).Preload("Dataset").First(&productDataset)
			result.Datasets = append(result.Datasets, productDataset.Dataset)
		}
	}

	renderx.JSON(w, http.StatusOK, result)
}
