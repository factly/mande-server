package product

import (
	"database/sql/driver"
	"errors"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
)

var headers = map[string]string{
	"X-User":         "1",
	"X-Organisation": "1",
}

var Product map[string]interface{} = map[string]interface{}{
	"title":              "Test Product",
	"slug":               "test-product",
	"price":              100,
	"status":             "teststatus",
	"currency_id":        1,
	"featured_medium_id": 1,
	"dataset_ids":        []uint{1},
	"tag_ids":            []uint{1},
}

var undecodableProduct map[string]interface{} = map[string]interface{}{
	"title":              4,
	"slug":               1,
	"price":              "100",
	"currency_id":        1,
	"featured_medium_id": 1,
	"dataset_ids":        "[]uint{1}",
	"tag_ids":            34,
}

var ProductReceive map[string]interface{} = map[string]interface{}{
	"title":              "Test Product",
	"slug":               "test-product",
	"price":              100,
	"status":             "teststatus",
	"currency_id":        1,
	"featured_medium_id": 1,
}

var invalidProduct map[string]interface{} = map[string]interface{}{
	"tie":                "Test Produ",
	"slg":                "test-prct",
	"price":              0,
	"status":             "tesatus",
	"currency_id":        1,
	"featured_medium_id": 1,
	"dataset_ids":        []uint{1},
	"tag_ids":            []uint{1},
}

var productlist []map[string]interface{} = []map[string]interface{}{
	{
		"title":              "Test Product 1",
		"slug":               "test-product-1",
		"price":              100,
		"status":             "teststatus1",
		"currency_id":        1,
		"featured_medium_id": 1,
		"dataset_ids":        []uint{1},
		"tag_ids":            []uint{1},
	},
	{
		"title":              "Test Product 2",
		"slug":               "test-product-2",
		"price":              200,
		"status":             "teststatus2",
		"currency_id":        1,
		"featured_medium_id": 1,
		"dataset_ids":        []uint{1},
		"tag_ids":            []uint{1},
	},
}

var ProductCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "title", "slug", "price", "status", "currency_id", "featured_medium_id"}
var catalogCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "featured_medium_id", "published_date"}
var planCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "name", "description", "status", "duration"}
var membershipCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "status", "user_id", "payment_id", "plan_id", "razorpay_order_id"}

var selectQuery string = `SELECT (.+) FROM "dp_product"`
var countQuery string = regexp.QuoteMeta(`SELECT count(1) FROM "dp_product"`)
var errProductMediumFK = errors.New(`pq: insert or update on table "dp_product" violates foreign key constraint "dp_product_featured_medium_id_dp_medium_id_foreign"`)
var errProductCurrencyFK = errors.New(`pq: insert or update on table "dp_product" violates foreign key constraint "dp_product_currency_id_dp_currency_id_foreign"`)

const basePath string = "/products"
const path string = "/products/{product_id}"

func ProductSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(selectQuery).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(ProductCols).
			AddRow(1, time.Now(), time.Now(), nil, Product["title"], Product["slug"], Product["price"], Product["status"], Product["currency_id"], Product["featured_medium_id"]))
}

func EmptyProductAssociationsMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_product_tag"`)).
		WillReturnRows(sqlmock.NewRows(append(tag.TagCols, []string{"tag_id", "product_id"}...)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_dataset" INNER JOIN "dp_product_dataset"`)).
		WillReturnRows(sqlmock.NewRows(append(dataset.DatasetCols, []string{"dataset_id", "product_id"}...)))
}

func tagsAssociationSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product_tag"`)).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "tag_id"}).
			AddRow(1, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows(tag.TagCols).
			AddRow(1, time.Now(), time.Now(), nil, tag.Tag["title"], tag.Tag["slug"]))
}

func datasetsAssociationSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product_dataset"`)).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "dataset_id"}).
			AddRow(1, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_dataset"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows(dataset.DatasetCols).
			AddRow(1, time.Now(), time.Now(), nil, dataset.Dataset["title"], dataset.Dataset["description"], dataset.Dataset["source"], dataset.Dataset["frequency"], dataset.Dataset["temporal_coverage"], dataset.Dataset["granularity"], dataset.Dataset["contact_name"], dataset.Dataset["contact_email"], dataset.Dataset["license"], dataset.Dataset["data_standard"], dataset.Dataset["sample_url"], dataset.Dataset["related_articles"], dataset.Dataset["time_saved"], dataset.Dataset["price"], dataset.Dataset["currency_id"], dataset.Dataset["featured_medium_id"]))
}

func catalogsAssociationSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_catalog" JOIN "dp_catalog_product"`)).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(append(catalogCols, []string{"product_id", "catalog_id"}...)).
			AddRow(1, time.Now(), time.Now(), nil, "title", "description", 1, time.Now(), 1, 1))
}

func plansCatalogAssociationSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(`SELECT (.+) FROM "dp_plan" INNER JOIN dp_plan_catalog`).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(append(planCols, []string{"catalog_id", "plan_id"}...)).
			AddRow(1, time.Now(), time.Now(), nil, "name", "description", "status", 10, 1, 1))
}

func membershipAssociationSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_membership"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows(membershipCols).
			AddRow(1, time.Now(), time.Now(), nil, "status", 1, 1, 1, "razorpay_order_id"))
}

func planSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_plan"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(planCols).
			AddRow(1, time.Now(), time.Now(), nil, "name", "description", "status", 10))
}

func insertWithErrorMock(mock sqlmock.Sqlmock, err error) {
	tag.TagSelectMock(mock)

	dataset.DatasetSelectMock(mock)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "dp_product"`).
		WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Product["title"], Product["slug"], Product["price"], Product["status"], Product["currency_id"], Product["featured_medium_id"]).
		WillReturnError(err)
	mock.ExpectRollback()
}

func preUpdateMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(ProductCols).
			AddRow(1, time.Now(), time.Now(), nil, "title", "slug", 200, "status", 2, 2))

	mock.ExpectBegin()

	tag.TagSelectMock(mock)

	mock.ExpectQuery(`INSERT INTO "dp_tag"`).
		WithArgs(test.AnyTime{}, test.AnyTime{}, nil, tag.Tag["title"], tag.Tag["slug"], 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec(`INSERT INTO "dp_product_tag"`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_tag"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	dataset.DatasetSelectMock(mock)

	mock.ExpectQuery(`INSERT INTO "dp_dataset"`).
		WithArgs(test.AnyTime{}, test.AnyTime{}, nil, dataset.Dataset["title"], dataset.Dataset["description"], dataset.Dataset["source"], dataset.Dataset["frequency"], dataset.Dataset["temporal_coverage"], dataset.Dataset["granularity"], dataset.Dataset["contact_name"], dataset.Dataset["contact_email"], dataset.Dataset["license"], dataset.Dataset["data_standard"], dataset.Dataset["sample_url"], dataset.Dataset["related_articles"], dataset.Dataset["time_saved"], dataset.Dataset["price"], dataset.Dataset["currency_id"], dataset.Dataset["featured_medium_id"], 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec(`INSERT INTO "dp_product_dataset"`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_dataset"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func updateMock(mock sqlmock.Sqlmock, err error) {
	preUpdateMock(mock)

	if err != nil {
		mock.ExpectExec(`UPDATE \"dp_product\"`).
			WithArgs(test.AnyTime{}, Product["title"], Product["slug"], Product["price"], Product["status"], Product["featured_medium_id"], Product["currency_id"], 1).
			WillReturnError(err)
		mock.ExpectRollback()
	} else {
		mock.ExpectExec(`UPDATE \"dp_product\"`).
			WithArgs(test.AnyTime{}, Product["title"], Product["slug"], Product["price"], Product["status"], Product["featured_medium_id"], Product["currency_id"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func updateMockWithoutMedium(mock sqlmock.Sqlmock) {
	preUpdateMock(mock)

	mock.ExpectExec(`UPDATE \"dp_product\"`).
		WithArgs(nil, test.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`UPDATE \"dp_product\"`).
		WithArgs(test.AnyTime{}, Product["title"], Product["slug"], Product["price"], Product["status"], Product["currency_id"], 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func productOrderExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "dp_order" JOIN`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func productCatalogExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "dp_catalog" JOIN`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func validateAssociations(result *httpexpect.Object) {
	result.Value("currency").
		Object().
		ContainsMap(currency.Currency)

	result.Value("featured_medium").
		Object().
		ContainsMap(medium.Medium)
}
