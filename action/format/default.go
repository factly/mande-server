package format

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
var DataFile = "./data/format.json"

// default - Create default format
// @Summary Create default format
// @Description Create default format
// @Tags Format
// @ID add-default-format
// @Consume json
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Success 201 {object} paging
// @Failure 400 {array} string
// @Router /formats/default [post]
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

	formats := make([]model.Format, 0)

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &formats)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx := model.DB.WithContext(context.WithValue(r.Context(), userContext, uID)).Begin()
	for i := range formats {
		tx.Model(&model.Format{}).FirstOrCreate(&formats[i], &formats[i])
	}

	result := paging{}
	result.Nodes = formats
	result.Total = int64(len(formats))

	tx.Commit()

	renderx.JSON(w, http.StatusCreated, result)
}
