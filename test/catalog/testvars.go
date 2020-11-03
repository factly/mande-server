package catalog

import (
	"database/sql/driver"
	"errors"
	"regexp"
	"time"

	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/tag"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/factly/data-portal-server/test/product"
	"github.com/gavv/httpexpect"
)

var Catalog map[string]interface{} = map[string]interface{}{
	"title":              "Test Title",
	"description":        "Test Description",
	"featured_medium_id": 1,
	"published_date":     time.Now(),
	"product_ids":        []uint{1},
}

var CatalogReceive map[string]interface{} = map[string]interface{}{
	"title":              "Test Title",
	"description":        "Test Description",
	"featured_medium_id": 1,
}

var invalidCatalog map[string]interface{} = map[string]interface{}{
	"tite":               "Test Tle",
	"descripon":          "Test Descri",
	"featured_medium_id": 1,
	"publisheddate":      nil,
}

var undecodableCatalog map[string]interface{} = map[string]interface{}{
	"title":              445,
	"description":        87,
	"featured_medium_id": "1",
}

var cataloglist []map[string]interface{} = []map[string]interface{}{
	{
		"title":              "Test Title 1",
		"description":        "Test Description 1",
		"featured_medium_id": 1,
		"published_date":     time.Now(),
		"product_ids":        []uint{1},
	},
	{
		"title":              "Test Title 2",
		"description":        "Test Description 2",
		"featured_medium_id": 1,
		"published_date":     time.Now(),
		"product_ids":        []uint{1},
	},
}

var CatalogCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "featured_medium_id", "published_date"}

var selectQuery string = `SELECT (.+) FROM \"dp_catalog\"`
var countQuery string = regexp.QuoteMeta(`SELECT count(1) FROM "dp_catalog"`)
var errCatalogProductFK = errors.New(`pq: insert or update on table "dp_catalog" violates foreign key constraint "dp_catalog_featured_medium_id_dp_medium_id_foreign"`)

const basePath string = "/catalogs"
const path string = "/catalogs/{catalog_id}"

func CatalogSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(selectQuery).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(CatalogCols).
			AddRow(1, time.Now(), time.Now(), nil, Catalog["title"], Catalog["description"], Catalog["featured_medium_id"], Catalog["published_date"]))
}

func productsAssociationSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_catalog_product"`)).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "catalog_id"}).
			AddRow(1, 1))
	product.ProductSelectMock(mock)
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

func preUpdateMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(CatalogCols).
			AddRow(1, time.Now(), time.Now(), nil, "title", "description", 1, time.Now()))

	mock.ExpectBegin()

	product.ProductSelectMock(mock)
}

func updateMock(mock sqlmock.Sqlmock, err error) {

	preUpdateMock(mock)

	mock.ExpectQuery(`INSERT INTO "dp_product"`).
		WithArgs(test.AnyTime{}, test.AnyTime{}, nil, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec(`INSERT INTO "dp_catalog_product"`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_catalog_product"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err != nil {
		mock.ExpectExec(`UPDATE \"dp_catalog\"`).
			WithArgs(test.AnyTime{}, Catalog["title"], Catalog["description"], Catalog["featured_medium_id"], test.AnyTime{}, 1).
			WillReturnError(err)
		mock.ExpectRollback()
	} else {
		mock.ExpectExec(`UPDATE \"dp_catalog\"`).
			WithArgs(test.AnyTime{}, Catalog["title"], Catalog["description"], Catalog["featured_medium_id"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func updateWithoutFeaturedMedium(mock sqlmock.Sqlmock) {
	preUpdateMock(mock)

	mock.ExpectQuery(`INSERT INTO "dp_product"`).
		WithArgs(test.AnyTime{}, test.AnyTime{}, nil, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec(`INSERT INTO "dp_catalog_product"`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_catalog_product"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`UPDATE \"dp_catalog\"`).
		WithArgs(nil, test.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	CatalogSelectMock(mock)

	mock.ExpectExec(`UPDATE \"dp_catalog\"`).
		WithArgs(test.AnyTime{}, Catalog["title"], Catalog["description"], test.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

}

func validateAssociations(result *httpexpect.Object) {
	result.Value("featured_medium").
		Object().
		ContainsMap(medium.Medium)
}
