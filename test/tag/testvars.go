package tag

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var Tag map[string]interface{} = map[string]interface{}{
	"title": "Test Tag",
	"slug":  "test-tag",
}

var invalidTag map[string]interface{} = map[string]interface{}{
	"titl": "Test",
	"slg":  "test",
}

var taglist []map[string]interface{} = []map[string]interface{}{
	{"title": "Test Tag 1", "slug": "test-tag-1"},
	{"title": "Test Tag 2", "slug": "test-tag-2"},
}

var TagCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_tag"`)

const basePath string = "/tags"
const path string = "/tags/{tag_id}"

func TagSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(TagCols).
			AddRow(1, time.Now(), time.Now(), nil, Tag["title"], Tag["slug"]))
}

func tagProductExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_product" INNER JOIN "dp_product_tag"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func tagDatasetExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_dataset" INNER JOIN "dp_dataset_tag"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}
