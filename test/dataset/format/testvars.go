package datasetformat

import (
	"errors"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var headers = map[string]string{
	"X-User":         "1",
	"X-Organisation": "1",
}

var DatasetFormat map[string]interface{} = map[string]interface{}{
	"format_id": 1,
	"url":       "test.url.com",
}

var datasetformatlist []map[string]interface{} = []map[string]interface{}{
	{
		"format_id":  1,
		"url":        "test1.url.com",
		"dataset_id": 1,
	},
	{
		"format_id":  1,
		"url":        "test2.url.com",
		"dataset_id": 1,
	},
}

var undecodableDatasetFormat map[string]interface{} = map[string]interface{}{
	"format_id": "1",
	"url":       15,
}

var invalidDatasetFormat map[string]interface{} = map[string]interface{}{
	"formatid": 1,
	"ur":       "test.url.com",
}

var DatasetFormatCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "created_by_id", "updated_by_id", "format_id", "dataset_id", "url"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_dataset_format"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(1) FROM "dp_dataset_format"`)
var errDatasetFormatFK error = errors.New(`pq: insert or update on table "dp_dataset_format" violates foreign key constraint "dp_dataset_format_format_id_dp_format_id_foreign"`)

const basePath string = "/datasets/{dataset_id}/format"
const path string = "/datasets/{dataset_id}/format/{format_id}"

func DatasetFormatSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(DatasetFormatCols).
			AddRow(1, time.Now(), time.Now(), nil, 1, 1, DatasetFormat["format_id"], 1, DatasetFormat["url"]))
}
