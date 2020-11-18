package plan

import (
	"database/sql/driver"
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

var selectQuery string = `SELECT (.+) FROM \"dp_plan\"`
var countQuery string = regexp.QuoteMeta(`SELECT count(1) FROM "dp_plan"`)

const basePath string = "/plans"
const path string = "/plans/{plan_id}"

func PlanSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(selectQuery).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(PlanCols).
			AddRow(1, time.Now(), time.Now(), nil, Plan["name"], Plan["description"], Plan["status"], Plan["duration"], Plan["price"], Plan["currency_id"]))
}

func associatedCatalogSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_plan_catalog"`)).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows([]string{"plan_id", "catalog_id"}).
			AddRow(1, 1))

	catalog.CatalogSelectMock(mock)
}

func productCatalogAssociationMock(mock sqlmock.Sqlmock, args ...driver.Value) {
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

func planMembershipExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "dp_membership`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func planUpdateMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(PlanCols).
			AddRow(1, time.Now(), time.Now(), nil, "name", "description", "status", 50, 100, 1))

	mock.ExpectBegin()
	catalog.CatalogSelectMock(mock)

	mock.ExpectQuery(`INSERT INTO "dp_catalog"`).
		WithArgs(test.AnyTime{}, test.AnyTime{}, nil, catalog.Catalog["title"], catalog.Catalog["description"], test.AnyTime{}, catalog.Catalog["featured_medium_id"], 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec(`INSERT INTO "dp_plan_catalog"`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_plan_catalog"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`UPDATE \"dp_plan\"`).
		WithArgs(test.AnyTime{}, Plan["name"], Plan["description"], Plan["price"], Plan["currency_id"], Plan["duration"], Plan["status"], 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	PlanSelectMock(mock)

	associatedCatalogSelectMock(mock)
	productCatalogAssociationMock(mock, 1)
	currency.CurrencySelectMock(mock)
	datasetsAssociationSelectMock(mock)
	tagsAssociationSelectMock(mock)
	currency.CurrencySelectMock(mock)
}
