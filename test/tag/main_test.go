package tag

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var tag map[string]interface{} = map[string]interface{}{
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

var tagCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_tag"`)

const basePath string = "/tags"
const path string = "/tags/{tag_id}"

func tagSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(tagCols).
			AddRow(1, time.Now(), time.Now(), nil, tag["title"], tag["slug"]))
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
