package user

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int          `json:"total"`
	Nodes []model.User `json:"nodes"`
}

// list - Get all users
// @Summary Show all users
// @Description Get all users
// @Tags User
// @ID get-all-users
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /users [get]
func list(w http.ResponseWriter, r *http.Request) {

	result := paging{}

	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Model(&model.User{}).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}
