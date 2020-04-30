package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// update - Update product by id
// @Summary Update a product by id
// @Description Update product by ID
// @Tags Product
// @ID update-product-by-id
// @Produce json
// @Consume json
// @Param id path string true "Product ID"
// @Param Product body product false "Product"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products/{id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	product := &model.Product{
		ID: uint(id),
	}

	req := &model.Product{}
	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&product).Updates(&model.Product{
		CurrencyID:    req.CurrencyID,
		ProductTypeID: req.ProductTypeID,
		StatusID:      req.StatusID,
		Title:         req.Title,
		Price:         req.Price,
		Slug:          req.Slug,
	})
	model.DB.First(&product).First(&product)
	model.DB.Model(&product).Association("ProductType").Find(&product.ProductType)
	model.DB.Model(&product).Association("Currency").Find(&product.Currency)
	model.DB.Model(&product).Association("Status").Find(&product.Status)

	json.NewEncoder(w).Encode(product)
}
