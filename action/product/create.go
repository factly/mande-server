package product

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-playground/validator/v10"
)

// createProduct - Create product
// @Summary Create product
// @Description Create product
// @Tags Product
// @ID add-product
// @Consume json
// @Produce  json
// @Param Product body product true "Product object"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products [post]
func createProduct(w http.ResponseWriter, r *http.Request) {

	req := &model.Product{}
	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.StructExcept(req, "ProductType", "Status", "Currency")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Product{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&req).Association("ProductType").Find(&req.ProductType)
	model.DB.Model(&req).Association("Currency").Find(&req.Currency)
	model.DB.Model(&req).Association("Status").Find(&req.Status)
	json.NewEncoder(w).Encode(req)
}
