package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// delete - Delete product by id
// @Summary Delete a product
// @Description Delete product by ID
// @Tags Product
// @ID delete-product-by-id
// @Consume  json
// @Param id path string true "Product ID"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products/{id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	product := &model.Product{
		ID: uint(id),
	}

	// check record exists or not
	err = model.DB.First(&product).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Model(&product).Association("ProductType").Find(&product.ProductType)
	model.DB.Model(&product).Association("Currency").Find(&product.Currency)
	model.DB.Model(&product).Association("Status").Find(&product.Status)
	model.DB.Delete(&product)

	json.NewEncoder(w).Encode(product)
}
