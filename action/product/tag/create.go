package tag

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

// create - create productTags
// @Summary Create productTags
// @Description create productTags
// @Tags ProductTag
// @ID add-productTags
// @Consume json
// @Produce  json
// @Param id path string true "Product ID"
// @Param ProductTag body productTags true "ProductTag object"
// @Success 201 {object} model.ProductTag
// @Failure 400 {array} string
// @Router /products/{id}/tag [post]
func create(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	productTag := &model.ProductTag{
		ProductID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&productTag)

	validate := validator.New()
	err = validate.Struct(productTag)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.ProductTag{}).Create(&productTag).Error

	if err != nil {
		log.Fatal(err)
	}

	render.JSON(w, http.StatusCreated, productTag)
}
