package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
)

// list - Get all products
// @Summary Show all products
// @Description Get all products
// @Tags Product
// @ID get-all-products
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Product
// @Router /products [get]
func list(w http.ResponseWriter, r *http.Request) {

	var products []model.Product
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

	model.DB.Offset(offset).Limit(limit).Preload("Currency").Preload("Status").Preload("ProductType").Model(&model.Product{}).Find(&products)

	json.NewEncoder(w).Encode(products)
}
