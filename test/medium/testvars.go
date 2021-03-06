package medium

import (
	"database/sql/driver"
	"regexp"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"

	"github.com/DATA-DOG/go-sqlmock"
)

var headers = map[string]string{
	"X-User":         "1",
	"X-Organisation": "1",
}

var Medium map[string]interface{} = map[string]interface{}{
	"name":        "Test Medium",
	"slug":        "test-medium",
	"type":        "testtype",
	"title":       "Test Title",
	"description": "Test Description",
	"caption":     "Test Caption",
	"alt_text":    "Test alt text",
	"file_size":   100,
	"url": postgres.Jsonb{
		RawMessage: []byte(`{"raw":"http://testimage.com/test.jpg"}`),
	},
	"dimensions": "testdims",
}

var undecodableMedium map[string]interface{} = map[string]interface{}{
	"name":      1,
	"slug":      3,
	"title":     99,
	"file_size": "100",
}

var invalidMedium map[string]interface{} = map[string]interface{}{
	"nam":         "Test Medium",
	"slug":        "test-medium",
	"type":        "testtype",
	"title":       "Test Title",
	"description": "Test Description",
	"caption":     "Test Caption",
	"alt_text":    "Test alt text",
	"filesize":    100,
	"url":         "http:/testurl.com",
	"dimensions":  "testdims",
}

var mediumlist []map[string]interface{} = []map[string]interface{}{
	{
		"name":        "Test Medium 1",
		"slug":        "test-medium-1",
		"type":        "testtype1",
		"title":       "Test Title 1",
		"description": "Test Description 1",
		"caption":     "Test Caption 1",
		"alt_text":    "Test alt text 1",
		"file_size":   100,
		"url": postgres.Jsonb{
			RawMessage: []byte(`{"raw":"http://testimage.com/test.jpg"}`),
		},
		"dimensions": "testdims1",
	},
	{
		"name":        "Test Medium 2",
		"slug":        "test-medium-2",
		"type":        "testtype2",
		"title":       "Test Title 2",
		"description": "Test Description 2",
		"caption":     "Test Caption 2",
		"alt_text":    "Test alt text 2",
		"file_size":   200,
		"url": postgres.Jsonb{
			RawMessage: []byte(`{"raw":"http://testimage.com/test.jpg"}`),
		},
		"dimensions": "testdims2",
	},
}

var MediumCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "created_by_id", "updated_by_id", "name", "slug", "type", "title", "description", "caption", "alt_text", "file_size", "url", "dimensions"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_medium"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(1) FROM "dp_medium"`)

const basePath string = "/media"
const path string = "/media/{media_id}"

func MediumSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(selectQuery).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(MediumCols).
			AddRow(1, time.Now(), time.Now(), nil, 1, 1, Medium["name"], Medium["slug"], Medium["type"], Medium["title"], Medium["description"], Medium["caption"], Medium["alt_text"], Medium["file_size"], Medium["url"], Medium["dimensions"]))

}

func mediumCatalogExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "dp_catalog`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func mediumDatasetExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "dp_dataset`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func mediumProductExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "dp_product`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}
