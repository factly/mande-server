package currency

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/factly/mande-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
)

// DataFile default json data file
var DataFile = "./data/currency.json"

// default - Create default currency
// @Summary Create default currency
// @Description Create default currency
// @Tags Currency
// @ID add-default-currency
// @Consume json
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Success 201 {object} paging
// @Failure 400 {array} string
// @Router /currencies/default [post]
func createDefault(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	jsonFile, err := os.Open(DataFile)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	defer jsonFile.Close()

	currencies := make([]model.Currency, 0)

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &currencies)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx := model.DB.WithContext(context.WithValue(r.Context(), userContext, uID)).Begin()
	for i := range currencies {
		tx.Model(&model.Currency{}).FirstOrCreate(&currencies[i], &currencies[i])
	}

	result := paging{}
	result.Nodes = currencies
	result.Total = int64(len(currencies))

	tx.Commit()

	renderx.JSON(w, http.StatusCreated, result)
}
