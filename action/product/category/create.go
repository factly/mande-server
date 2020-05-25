package category

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// create - create productCategory
// @Summary Create productCategory
// @Description create productCategory
// @Tags ProductCategory
// @ID add-productCategory
// @Consume json
// @Produce  json
// @Param product_id path string true "Product ID"
// @Param ProductCategory body productCategory true "ProductCategory object"
// @Success 201 {object} model.ProductCategory
// @Failure 400 {array} string
// @Router /products/{product_id}/category [post]
func create(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	productCategory := &productCategory{}
	result := &model.ProductCategory{}

	json.NewDecoder(r.Body).Decode(&productCategory)

	validate := validator.New()
	err = validate.Struct(productCategory)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	result.ProductID = uint(id)
	result.CategoryID = productCategory.CategoryID

	err = model.DB.Model(&model.ProductCategory{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}

	render.JSON(w, http.StatusCreated, result)
}
