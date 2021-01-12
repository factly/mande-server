package product

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list - Get all products owned by user
// @Summary Show all products owned by user
// @Description Get all products owned by user
// @Tags Product
// @ID get-all-products-owned-by-user
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /products/my [get]
func my(w http.ResponseWriter, r *http.Request) {

	uID, err := util.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &paging{}
	result.Nodes = make([]model.Product, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	orders := []model.Order{}

	model.DB.Preload("Products").Preload("Products.Datasets").Preload("Products.Tags").Preload("Products.Currency").Preload("Products.FeaturedMedium").Model(&model.Order{}).Where(&model.Order{
		UserID: uint(uID),
	}).Find(&orders)

	products := []model.Product{}

	for _, order := range orders {
		products = append(products, order.Products...)
	}

	result.Nodes = products[offset : offset+limit]
	result.Total = int64(len(products))

	renderx.JSON(w, http.StatusOK, result)
}
