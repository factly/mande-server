package category

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
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
// @Param id path string true "Product ID"
// @Param ProductCategory body productCategory true "ProductCategory object"
// @Success 200 {object} model.ProductCategory
// @Failure 400 {array} string
// @Router /products/{id}/category [post]
func create(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.ProductCategory{
		ProductID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}
	err = model.DB.Model(&model.ProductCategory{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}
