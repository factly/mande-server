package product

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
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
		renderx.JSON(w, http.StatusBadRequest, validationError)
		return
	}

	result := &productData{}
	result.Product = model.Product{
		Title:      product.Title,
		Slug:       product.Slug,
		Price:      product.Price,
		Status:     product.Status,
		CurrencyID: product.CurrencyID,
	}

	err := model.DB.Model(&model.Product{}).Create(&result.Product).Error

	if err != nil {
		log.Fatal(err)
	}

	model.DB.Preload("Currency").First(&result.Product)

	for _, id := range product.DatasetIDs {
		productDataset := &model.ProductDataset{}

		productDataset.DatasetID = uint(id)
		productDataset.ProductID = result.ID

		err = model.DB.Model(&model.ProductDataset{}).Create(&productDataset).Error

		if err != nil {
			log.Fatal(err)
		}
		model.DB.Model(&model.ProductDataset{}).Preload("Dataset").First(&productDataset)
		result.Datasets = append(result.Datasets, productDataset.Dataset)
	}

	for _, id := range product.TagIDs {
		productTag := &model.ProductTag{}

		productTag.TagID = uint(id)
		productTag.ProductID = result.ID

		err = model.DB.Model(&model.ProductTag{}).Create(&productTag).Error

		if err != nil {
			log.Fatal(err)
		}
		model.DB.Model(&model.ProductTag{}).Preload("Tag").First(&productTag)
		result.Tags = append(result.Tags, productTag.Tag)
	}

	renderx.JSON(w, http.StatusCreated, result)
}
