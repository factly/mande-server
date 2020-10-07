package plan

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/catalog"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/product"
	"github.com/factly/data-portal-server/test/tag"
)

var Plan map[string]interface{} = map[string]interface{}{
	"name":        "Test Plan",
	"description": "Test Plan Description",
	"status":      "teststatus",
	"price":       100,
	"currency_id": 1,
	"duration":    364,
	"catalog_ids": []uint{1},
}

var PlanReceive map[string]interface{} = map[string]interface{}{
	"name":        "Test Plan",
	"description": "Test Plan Description",
	"status":      "teststatus",
	"price":       100,
	"currency_id": 1,
	"duration":    364,
}

var undecodablePlan map[string]interface{} = map[string]interface{}{
	"name":        "Test Plan",
	"description": 4322,
	"status":      "teststatus",
	"duration":    "364",
}

var invalidPlan map[string]interface{} = map[string]interface{}{
	"planname":  "Test Plan",
	"plan_info": "Test Plan Info",
	"status":    "teststatus",
}

var planlist []map[string]interface{} = []map[string]interface{}{
	{
		"name":        "Test Plan 1",
		"description": "Test Plan Description 1",
		"status":      "teststatus1",
		"price":       100,
		"currency_id": 1,
		"duration":    364,
		"catalog_ids": []uint{1},
	},
	{
		"name":        "Test Plan 2",
		"description": "Test Plan Description 2",
		"status":      "teststatus2",
		"price":       200,
		"currency_id": 1,
		"duration":    223,
		"catalog_ids": []uint{1},
	},
}

var PlanCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "name", "description", "status", "duration", "price", "currency_id"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_plan"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_plan"`)

const basePath string = "/plans"
const path string = "/plans/{plan_id}"

func PlanSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(PlanCols).
			AddRow(1, time.Now(), time.Now(), nil, Plan["name"], Plan["description"], Plan["status"], Plan["duration"], Plan["price"], Plan["currency_id"]))
}

func associatedCatalogSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_catalog" INNER JOIN "dp_plan_catalog"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows(append(catalog.CatalogCols, []string{"plan_id", "catalog_id"}...)).
			AddRow(1, time.Now(), time.Now(), nil, catalog.Catalog["title"], catalog.Catalog["description"], catalog.Catalog["featured_medium_id"], catalog.Catalog["published_date"], 1, 1))
}

func associatedCatalogWithArg(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_catalog" INNER JOIN "dp_plan_catalog"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows(append(catalog.CatalogCols, []string{"plan_id", "catalog_id"}...)).
			AddRow(1, time.Now(), time.Now(), nil, catalog.Catalog["title"], catalog.Catalog["description"], catalog.Catalog["featured_medium_id"], catalog.Catalog["published_date"], 1, 1))
}

func productCatalogAssociationMock(mock sqlmock.Sqlmock, catId uint) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product" INNER JOIN "dp_catalog_product"`)).
		WithArgs(catId).
		WillReturnRows(sqlmock.NewRows(append(product.ProductCols, []string{"product_id", "catalog_id"}...)).
			AddRow(1, time.Now(), time.Now(), nil, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1, catId))
}

func planMembershipExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_membership`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func planUpdateMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(PlanCols).
			AddRow(1, time.Now(), time.Now(), nil, "name", "description", "status", 50, 100, 1))

	associatedCatalogSelectMock(mock)

	mock.ExpectBegin()
	catalog.CatalogSelectMock(mock)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_plan_catalog"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`UPDATE \"dp_plan\" SET (.+)  WHERE (.+) \"dp_plan\".\"id\" = `).
		WithArgs(Plan["currency_id"], Plan["description"], Plan["duration"], Plan["name"], Plan["price"], Plan["status"], test.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO "dp_plan_catalog"`).
		WithArgs(1, 1, 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	PlanSelectMock(mock)

	currency.CurrencySelectMock(mock)

	associatedCatalogSelectMock(mock)

	productCatalogAssociationMock(mock, 1)

	currency.CurrencySelectMock(mock)

	dataset.DatasetSelectMock(mock)

	tag.TagSelectMock(mock)

}
