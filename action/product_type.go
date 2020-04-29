package action

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-playground/validator/v10"
)

// productType request body
type productType struct {
	Name string `json:"name"`
}

// GetProductTypes - Get all productTypes
// @Summary Show all productTypes
// @Description Get all productTypes
// @Tags Type
// @ID get-all-productTypes
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.ProductType
// @Router /products/{id}/type [get]
func GetProductTypes(w http.ResponseWriter, r *http.Request) {

	var productTypes []model.ProductType
	p := r.URL.Query().Get("page")
	pg, _ := strconv.Atoi(p) // pg contains page number
	l := r.URL.Query().Get("limit")
	li, _ := strconv.Atoi(l) // li contains perPage number

	offset := 0 // no. of records to skip
	limit := 5  // limt

	if li > 0 && li <= 10 {
		limit = li
	}

	if pg > 1 {
		offset = (pg - 1) * limit
	}

	model.DB.Offset(offset).Limit(limit).Model(&model.ProductType{}).Find(&productTypes)

	json.NewEncoder(w).Encode(productTypes)
}

// CreateProductType - Create product type
// @Summary Create product type
// @Description Create product type
// @Tags Type
// @ID add-type
// @Consume json
// @Produce  json
// @Param Type body productType true "Type object"
// @Success 200 {object} model.ProductType
// @Failure 400 {array} string
// @Router /products/{id}/type [post]
func CreateProductType(w http.ResponseWriter, r *http.Request) {

	req := &model.ProductType{}
	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.ProductType{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(req)
}
