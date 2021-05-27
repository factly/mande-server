package format

import (
	"database/sql/driver"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var headers = map[string]string{
	"X-User":         "1",
	"X-Organisation": "1",
}

var Format map[string]interface{} = map[string]interface{}{
	"name":        "Test Format",
	"description": "Test Description",
	"is_default":  true,
}

var undecodableFormat map[string]interface{} = map[string]interface{}{
	"name":        10,
	"description": 20,
	"is_default":  "true",
}

var invalidFormat map[string]interface{} = map[string]interface{}{
	"nae":        "Test Format",
	"decription": "Test Description",
	"isdefault":  true,
}

var formatlist []map[string]interface{} = []map[string]interface{}{
	{
		"name":        "Test Format 1",
		"description": "Test Description 1",
		"is_default":  true,
	},
	{
		"name":        "Test Format 2",
		"description": "Test Description 2",
		"is_default":  false,
	},
}

var defaultFormats []map[string]interface{} = []map[string]interface{}{
	{
		"name":        "JSON",
		"description": "JSON",
		"is_default":  false,
	},
	{
		"name":        "CSV",
		"description": "CSV",
		"is_default":  false,
	},
	{
		"name":        "PDF",
		"description": "PDF",
		"is_default":  false,
	},
}

var FormatCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "created_by_id", "updated_by_id", "name", "description", "is_default"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_format"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(1) FROM "dp_format"`)

const basePath string = "/formats"
const path string = "/formats/{format_id}"
const defaultPath string = "/formats/default"

func FormatSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(selectQuery).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(FormatCols).
			AddRow(1, time.Now(), time.Now(), nil, 1, 1, Format["name"], Format["description"], Format["is_default"]))
}

func formatDatasetExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "dp_dataset_format"  WHERE`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}
