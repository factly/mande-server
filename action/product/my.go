package product

import (
	"net/http"

	"github.com/factly/mande-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// userDetails - Get product by id
// @Summary Show a product by id
// @Description Get product by ID
// @Tags Product
// @ID get-my-products
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} productRes
// @Failure 400 {array} string
// @Router /products/my [get]
func my(w http.ResponseWriter, r *http.Request) {

	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	offset, limit := paginationx.Parse(r.URL.Query())

	result := &paging{}

	orders := []model.Order{}

	products := make([]model.Product, 0)

	model.DB.Preload("Products").Model(&model.Order{}).Where(&model.Order{
		Status: "complete",
		UserID: uint(uID),
	}).Find(&orders)

	for _, each := range orders {
		products = append(products, each.Products...)
	}

	result.Total = int64(len(products))
	result.Nodes = products[offset : offset+limit]

	renderx.JSON(w, http.StatusOK, result)
}
