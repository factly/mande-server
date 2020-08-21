package format

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var format map[string]interface{} = map[string]interface{}{
	"name":        "Test Format",
	"description": "Test Description",
	"is_default":  true,
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

var formatCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "name", "description", "is_default"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_format"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_format"`)

const basePath string = "/formats"
const path string = "/formats/{format_id}"

func formatSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(formatCols).
			AddRow(1, time.Now(), time.Now(), nil, format["name"], format["description"], format["is_default"]))
}

func formatDatasetExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_dataset_format"  WHERE`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}
