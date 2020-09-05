package product

import (
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

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_product"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_product"`)
var errProductMediumFK = errors.New(`pq: insert or update on table "dp_product" violates foreign key constraint "dp_product_featured_medium_id_dp_medium_id_foreign"`)
var errProductCurrencyFK = errors.New(`pq: insert or update on table "dp_product" violates foreign key constraint "dp_product_currency_id_dp_currency_id_foreign"`)

const basePath string = "/products"
const path string = "/products/{product_id}"

func ProductSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(ProductCols).
			AddRow(1, time.Now(), time.Now(), nil, Product["title"], Product["slug"], Product["price"], Product["status"], Product["currency_id"], Product["featured_medium_id"]))
}

func EmptyProductAssociationsMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_product_tag"`)).
		WillReturnRows(sqlmock.NewRows(append(tag.TagCols, []string{"tag_id", "product_id"}...)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_dataset" INNER JOIN "dp_product_dataset"`)).
		WillReturnRows(sqlmock.NewRows(append(dataset.DatasetCols, []string{"dataset_id", "product_id"}...)))
}

func tagsAssociationSelectMock(mock sqlmock.Sqlmock, prodId int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_product_tag"`)).
		WithArgs(prodId).
		WillReturnRows(sqlmock.NewRows(append(tag.TagCols, []string{"tag_id", "product_id"}...)).
			AddRow(1, time.Now(), time.Now(), nil, tag.Tag["title"], tag.Tag["slug"], 1, prodId))
}

func datasetsAssociationSelectMock(mock sqlmock.Sqlmock, prodId int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_dataset" INNER JOIN "dp_product_dataset"`)).
		WithArgs(prodId).
		WillReturnRows(sqlmock.NewRows(append(dataset.DatasetCols, []string{"dataset_id", "product_id"}...)).
			AddRow(1, time.Now(), time.Now(), nil, dataset.Dataset["title"], dataset.Dataset["description"], dataset.Dataset["source"], dataset.Dataset["frequency"], dataset.Dataset["temporal_coverage"], dataset.Dataset["granularity"], dataset.Dataset["contact_name"], dataset.Dataset["contact_email"], dataset.Dataset["license"], dataset.Dataset["data_standard"], dataset.Dataset["related_articles"], dataset.Dataset["time_saved"], dataset.Dataset["price"], dataset.Dataset["currency_id"], dataset.Dataset["featured_medium_id"], 1, prodId))
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

	tagsAssociationSelectMock(mock, 1)

	datasetsAssociationSelectMock(mock, 1)

	mock.ExpectBegin()

	tag.TagSelectMock(mock)

	dataset.DatasetSelectMock(mock)

}

func updateMock(mock sqlmock.Sqlmock, err error) {
	preUpdateMock(mock)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_tag"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_dataset"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err != nil {
		mock.ExpectExec(`UPDATE \"dp_product\" SET (.+)  WHERE (.+) \"dp_product\".\"id\" = `).
			WithArgs(Product["currency_id"], Product["featured_medium_id"], Product["price"], Product["slug"], Product["status"], Product["title"], test.AnyTime{}, 1).
			WillReturnError(err)
		mock.ExpectRollback()
	} else {
		mock.ExpectExec(`UPDATE \"dp_product\" SET (.+)  WHERE (.+) \"dp_product\".\"id\" = `).
			WithArgs(Product["currency_id"], Product["featured_medium_id"], Product["price"], Product["slug"], Product["status"], Product["title"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`INSERT INTO "dp_product_tag"`).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(`INSERT INTO "dp_product_dataset"`).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
}

func updateMockWithoutMedium(mock sqlmock.Sqlmock) {
	preUpdateMock(mock)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_tag"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_dataset"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`UPDATE \"dp_product\" SET (.+)  WHERE (.+) \"dp_product\".\"id\" = `).
		WithArgs(nil, test.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ProductSelectMock(mock)

	mock.ExpectExec(`UPDATE \"dp_product\" SET (.+)  WHERE (.+) \"dp_product\".\"id\" = `).
		WithArgs(Product["currency_id"], Product["price"], Product["slug"], Product["status"], Product["title"], test.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO "dp_product_tag"`).
		WithArgs(1, 1, 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`INSERT INTO "dp_product_dataset"`).
		WithArgs(1, 1, 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

}

func productCartExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_cart" INNER JOIN`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func productCatalogExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_catalog" INNER JOIN`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func productOrderExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_order_item"`)).
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
