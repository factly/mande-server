package cart

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int64            `json:"total"`
	Nodes []model.CartItem `json:"nodes"`
}

// list - Get all carts
// @Summary Show all carts
// @Description Get all carts
// @Tags Cart
// @ID get-all-carts
// @Produce  json
// @Param X-User header string false "User ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /cartitems [get]
func list(w http.ResponseWriter, r *http.Request) {
	uID, err := util.GetUser(r)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := paging{}
	result.Nodes = make([]model.CartItem, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Model(&model.CartItem{}).Where(&model.CartItem{
		UserID: uint(uID),
	}).Preload("Product").Preload("Product.Currency").Preload("Product.FeaturedMedium").Preload("Product.Tags").Preload("Product.Datasets").Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}
