package dataset

import (
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/format"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
	"github.com/jinzhu/gorm/dialects/postgres"
)

func nilJsonb() postgres.Jsonb {
	ba, _ := json.Marshal(nil)
	return postgres.Jsonb{
		RawMessage: ba,
	}
}

var Dataset map[string]interface{} = map[string]interface{}{
	"title":              "Test Title",
	"description":        "Test Description",
	"source":             "testsource",
	"frequency":          "testfreq",
	"temporal_coverage":  "Test coverage",
	"granularity":        "test",
	"contact_name":       "Test Name",
	"contact_email":      "test@mail.com",
	"license":            "TestLicense",
	"data_standard":      "Test Datastd",
	"sample_url":         "testurl.com",
	"related_articles":   nilJsonb(),
	"time_saved":         10,
	"price":              100,
	"currency_id":        1,
	"featured_medium_id": 1,
	"tag_ids":            []uint{1},
}

var DatasetReceive map[string]interface{} = map[string]interface{}{
	"title":              "Test Title",
	"description":        "Test Description",
	"source":             "testsource",
	"frequency":          "testfreq",
	"temporal_coverage":  "Test coverage",
	"granularity":        "test",
	"contact_name":       "Test Name",
	"contact_email":      "test@mail.com",
	"license":            "TestLicense",
	"data_standard":      "Test Datastd",
	"sample_url":         "testurl.com",
	"related_articles":   nil,
	"time_saved":         10,
	"price":              100,
	"currency_id":        1,
	"featured_medium_id": 1,
}

var invalidDataset map[string]interface{} = map[string]interface{}{
	"tite":               "Test Titl",
	"desciption":         "Test Desc",
	"source":             "testsource",
	"temporal_coverage":  "Test cov",
	"granularity":        "test",
	"pric":               100,
	"currency_id":        1,
	"featured_medium_id": 1,
	"tag_ids":            []uint{1},
}

var undecodableDataset map[string]interface{} = map[string]interface{}{
	"tite":               23,
	"desciption":         42,
	"featured_medium_id": "1",
}

var datasetlist []map[string]interface{} = []map[string]interface{}{
	{
		"title":              "Test Title 1",
		"description":        "Test Description 1",
		"source":             "testsource1",
		"frequency":          "testfreq1",
		"temporal_coverage":  "Test coverage 1",
		"granularity":        "test1",
		"contact_name":       "Test Name 1",
		"contact_email":      "test1@mail.com",
		"license":            "TestLicense1",
		"data_standard":      "Test Datastd 1",
		"sample_url":         "testurl1.com",
		"related_articles":   nilJsonb(),
		"time_saved":         10,
		"price":              100,
		"currency_id":        1,
		"featured_medium_id": 1,
		"tag_ids":            []uint{1},
	},
	{
		"title":              "Test Title 2",
		"description":        "Test Description 2",
		"source":             "testsource2",
		"frequency":          "testfreq2",
		"temporal_coverage":  "Test coverage 2",
		"granularity":        "test2",
		"contact_name":       "Test Name 2",
		"contact_email":      "test2@mail.com",
		"license":            "TestLicense2",
		"data_standard":      "Test Datastd 2",
		"sample_url":         "testurl2.com",
		"related_articles":   nilJsonb(),
		"time_saved":         20,
		"price":              200,
		"currency_id":        1,
		"featured_medium_id": 1,
		"tag_ids":            []uint{1},
	},
}

var DatasetCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "source", "frequency", "temporal_coverage", "granularity", "contact_name", "contact_email", "license", "data_standard", "sample_url", "related_articles", "time_saved", "price", "currency_id", "featured_medium_id"}
var productCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "title", "slug", "price", "status", "currency_id", "featured_medium_id"}
var orderCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "user_id", "status", "payment_id", "razorpay_order_id"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_dataset"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_dataset"`)
var errDatasetMediumFK = errors.New(`pq: insert or update on table "dp_dataset" violates foreign key constraint "dp_dataset_featured_medium_id_dp_medium_id_foreign"`)
var errDatasetCurrencyFK = errors.New(`pq: insert or update on table "dp_dataset" violates foreign key constraint "dp_dataset_currency_id_dp_currency_id_foreign"`)

const basePath string = "/datasets"
const path string = "/datasets/{dataset_id}"

func DatasetSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(DatasetCols).
			AddRow(1, time.Now(), time.Now(), nil, Dataset["title"], Dataset["description"], Dataset["source"], Dataset["frequency"], Dataset["temporal_coverage"], Dataset["granularity"], Dataset["contact_name"], Dataset["contact_email"], Dataset["license"], Dataset["data_standard"], Dataset["sample_url"], Dataset["related_articles"], Dataset["time_saved"], Dataset["price"], Dataset["currency_id"], Dataset["featured_medium_id"]))
}
func orderSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(orderCols).
			AddRow(1, time.Now(), time.Now(), nil, 1, "status", 1, "razorpay_order_id"))
}
func tagAssociationSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_dataset_tag"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(append(tag.TagCols, []string{"tag_id", "dataset_id"}...)).
			AddRow(1, time.Now(), time.Now(), nil, "title1", "slug1", 1, 1))
}

func datasetFormatSelectMock(mock sqlmock.Sqlmock, id int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_dataset_format"`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "format_id", "dataset_id", "url"}).
			AddRow(id, time.Now(), time.Now(), nil, 1, id, "www.testurl.com"))

	format.FormatSelectMock(mock)
}

func validateAssociations(result *httpexpect.Object) {
	result.Value("featured_medium").
		Object().
		ContainsMap(medium.Medium)

	result.Value("currency").
		Object().
		ContainsMap(currency.Currency)
}

func insertWithErrorMock(mock sqlmock.Sqlmock, err error) {
	tag.TagSelectMock(mock)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "dp_dataset"`).
		WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Dataset["title"], Dataset["description"], Dataset["source"], Dataset["frequency"], Dataset["temporal_coverage"], Dataset["granularity"], Dataset["contact_name"], Dataset["contact_email"], Dataset["license"], Dataset["data_standard"], Dataset["sample_url"], Dataset["related_articles"], Dataset["time_saved"], Dataset["price"], Dataset["currency_id"], Dataset["featured_medium_id"]).
		WillReturnError(err)
	mock.ExpectRollback()
}

func preUpdateMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(DatasetCols).
			AddRow(1, time.Now(), time.Now(), nil, "title", "description", "source", "frequency", "temporal_coverage", "granularity", "contact_name", "contact_email", "license", "data_standard", "sample_url", nilJsonb(), 10, 200, 2, 2))

	tagAssociationSelectMock(mock)

	mock.ExpectBegin()

	tag.TagSelectMock(mock)
}

func updateMock(mock sqlmock.Sqlmock, err error) {
	preUpdateMock(mock)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_dataset_tag"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err != nil {
		mock.ExpectExec(`UPDATE \"dp_dataset\" SET (.+)  WHERE (.+) \"dp_dataset\".\"id\" = `).
			WithArgs(Dataset["contact_email"], Dataset["contact_name"], Dataset["currency_id"], Dataset["data_standard"], Dataset["description"], Dataset["featured_medium_id"], Dataset["frequency"], Dataset["granularity"], Dataset["license"], Dataset["price"], Dataset["related_articles"], Dataset["sample_url"], Dataset["source"], Dataset["temporal_coverage"], Dataset["time_saved"], Dataset["title"], test.AnyTime{}, 1).
			WillReturnError(err)
		mock.ExpectRollback()
	} else {
		mock.ExpectExec(`UPDATE \"dp_dataset\" SET (.+)  WHERE (.+) \"dp_dataset\".\"id\" = `).
			WithArgs(Dataset["contact_email"], Dataset["contact_name"], Dataset["currency_id"], Dataset["data_standard"], Dataset["description"], Dataset["featured_medium_id"], Dataset["frequency"], Dataset["granularity"], Dataset["license"], Dataset["price"], Dataset["related_articles"], Dataset["sample_url"], Dataset["source"], Dataset["temporal_coverage"], Dataset["time_saved"], Dataset["title"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`INSERT INTO "dp_dataset_tag"`).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
}

func updateWithoutFeaturedMedium(mock sqlmock.Sqlmock) {
	preUpdateMock(mock)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_dataset_tag"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`UPDATE \"dp_dataset\" SET (.+)  WHERE (.+) \"dp_dataset\".\"id\" = `).
		WithArgs(nil, test.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	DatasetSelectMock(mock)

	mock.ExpectExec(`UPDATE \"dp_dataset\" SET (.+)  WHERE (.+) \"dp_dataset\".\"id\" = `).
		WithArgs(Dataset["contact_email"], Dataset["contact_name"], Dataset["currency_id"], Dataset["data_standard"], Dataset["description"], Dataset["frequency"], Dataset["granularity"], Dataset["license"], Dataset["price"], Dataset["related_articles"], Dataset["sample_url"], Dataset["source"], Dataset["temporal_coverage"], Dataset["time_saved"], Dataset["title"], test.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO "dp_dataset_tag"`).
		WithArgs(1, 1, 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func datasetProductExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_product" INNER JOIN`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func deleteMock(mock sqlmock.Sqlmock, err error) {
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_dataset_format" SET "deleted_at"=`)).
		WithArgs(test.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_dataset_tag"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err == nil {
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_dataset" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
	} else {
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_dataset" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnError(err)
	}
}
