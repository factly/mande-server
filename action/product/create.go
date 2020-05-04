package product

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-playground/validator/v10"
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

	product := &model.Product{}
	json.NewDecoder(r.Body).Decode(&product)

	validate := validator.New()
	err := validate.StructExcept(product, "ProductType", "Status", "Currency")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Product{}).Create(&product).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&product).Association("ProductType").Find(&product.ProductType)
	model.DB.Model(&product).Association("Currency").Find(&product.Currency)
	model.DB.Model(&product).Association("Status").Find(&product.Status)

	render.JSON(w, http.StatusCreated, product)
}
