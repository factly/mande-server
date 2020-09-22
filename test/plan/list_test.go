package plan

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/catalog"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/product"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
)

func TestListPlan(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get empty plan list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(PlanCols))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_catalog" INNER JOIN "dp_plan_catalog"`)).
			WillReturnRows(sqlmock.NewRows(append(catalog.CatalogCols, []string{"plan_id", "catalog_id"}...)))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product" INNER JOIN "dp_catalog_product"`)).
			WillReturnRows(sqlmock.NewRows(append(product.ProductCols, []string{"product_id", "catalog_id"}...)))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_dataset" INNER JOIN "dp_product_dataset"`)).
			WillReturnRows(sqlmock.NewRows(append(dataset.DatasetCols, []string{"dataset_id", "product_id"}...)))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_product_tag"`)).
			WillReturnRows(sqlmock.NewRows(append(tag.TagCols, []string{"tag_id", "product_id"}...)))

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("list plans", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(planlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(PlanCols).
				AddRow(1, time.Now(), time.Now(), nil, planlist[0]["name"], planlist[0]["description"], planlist[0]["status"], planlist[0]["duration"]).
				AddRow(2, time.Now(), time.Now(), nil, planlist[1]["name"], planlist[1]["description"], planlist[1]["status"], planlist[1]["duration"]))

		associatedCatalogWithArg(mock)

		productCatalogAssociationMock(mock, 1)

		currency.CurrencySelectMock(mock)

		dataset.DatasetSelectMock(mock)

		tag.TagSelectMock(mock)

		delete(planlist[0], "catalog_ids")
		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(planlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(planlist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("get plan list with paiganation", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(planlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(PlanCols).
				AddRow(2, time.Now(), time.Now(), nil, planlist[1]["name"], planlist[1]["description"], planlist[1]["status"], planlist[1]["duration"]))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_catalog" INNER JOIN "dp_plan_catalog"`)).
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows(append(catalog.CatalogCols, []string{"plan_id", "catalog_id"}...)).
				AddRow(1, time.Now(), time.Now(), nil, catalog.Catalog["title"], catalog.Catalog["description"], catalog.Catalog["featured_medium_id"], catalog.Catalog["published_date"], 1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product" INNER JOIN "dp_catalog_product"`)).
			WillReturnRows(sqlmock.NewRows(append(product.ProductCols, []string{"product_id", "catalog_id"}...)))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_dataset" INNER JOIN "dp_product_dataset"`)).
			WillReturnRows(sqlmock.NewRows(append(dataset.DatasetCols, []string{"dataset_id", "product_id"}...)))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_product_tag"`)).
			WillReturnRows(sqlmock.NewRows(append(tag.TagCols, []string{"tag_id", "product_id"}...)))

		delete(planlist[1], "catalog_ids")
		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(planlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(planlist[1])

		test.ExpectationsMet(t, mock)
	})
}
