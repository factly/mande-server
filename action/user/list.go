package user

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list - Get all users
// @Summary Show all users
// @Description Get all users
// @Tags User
// @ID get-all-users
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.User
// @Router /users [get]
func list(w http.ResponseWriter, r *http.Request) {

	var users []model.User

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Offset(offset).Limit(limit).Model(&model.User{}).Find(&users)

	json.NewEncoder(w).Encode(users)
}
