package catalog

import (
	"errors"
	"regexp"
	"time"

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

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_catalog"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_catalog"`)
var errCatalogProductFK = errors.New(`pq: insert or update on table "dp_catalog" violates foreign key constraint "dp_catalog_featured_medium_id_dp_medium_id_foreign"`)

const basePath string = "/catalogs"
const path string = "/catalogs/{catalog_id}"

func CatalogSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(CatalogCols).
			AddRow(1, time.Now(), time.Now(), nil, Catalog["title"], Catalog["description"], Catalog["featured_medium_id"], Catalog["published_date"]))
}

func productsAssociationSelectMock(mock sqlmock.Sqlmock, catId int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product" INNER JOIN "dp_catalog_product"`)).
		WithArgs(catId).
		WillReturnRows(sqlmock.NewRows(append(product.ProductCols, []string{"product_id", "catalog_id"}...)).
			AddRow(1, time.Now(), time.Now(), nil, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1, catId))
}

func preUpdateMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(CatalogCols).
			AddRow(1, time.Now(), time.Now(), nil, "title", "description", 1, time.Now()))

	productsAssociationSelectMock(mock, 1)

	mock.ExpectBegin()

	product.ProductSelectMock(mock)

}

func updateMock(mock sqlmock.Sqlmock, err error) {

	preUpdateMock(mock)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_catalog_product"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err != nil {
		mock.ExpectExec(`UPDATE \"dp_catalog\" SET (.+)  WHERE (.+) \"dp_catalog\".\"id\" = `).
			WithArgs(Catalog["description"], Catalog["featured_medium_id"], test.AnyTime{}, Catalog["title"], test.AnyTime{}, 1).
			WillReturnError(err)
		mock.ExpectRollback()
	} else {
		mock.ExpectExec(`UPDATE \"dp_catalog\" SET (.+)  WHERE (.+) \"dp_catalog\".\"id\" = `).
			WithArgs(Catalog["description"], Catalog["featured_medium_id"], test.AnyTime{}, Catalog["title"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`INSERT INTO "dp_catalog_product"`).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
}

func updateWithoutFeaturedMedium(mock sqlmock.Sqlmock) {
	preUpdateMock(mock)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_catalog_product"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`UPDATE \"dp_catalog\" SET (.+)  WHERE (.+) \"dp_catalog\".\"id\" = `).
		WithArgs(nil, test.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	CatalogSelectMock(mock)

	mock.ExpectExec(`UPDATE \"dp_catalog\" SET (.+)  WHERE (.+) \"dp_catalog\".\"id\" = `).
		WithArgs(Catalog["description"], test.AnyTime{}, Catalog["title"], test.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO "dp_catalog_product"`).
		WithArgs(1, 1, 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

}

func validateAssociations(result *httpexpect.Object) {
	result.Value("featured_medium").
		Object().
		ContainsMap(medium.Medium)
}
