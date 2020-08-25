package catalog

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/factly/data-portal-server/test/product"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
)

func TestListCatalog(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get empty catalog list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(CatalogCols))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product" INNER JOIN "dp_catalog_product"`)).
			WillReturnRows(sqlmock.NewRows(append(product.ProductCols, []string{"product_id", "catalog_id"}...)))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_product_tag"`)).
			WillReturnRows(sqlmock.NewRows(append(tag.TagCols, []string{"tag_id", "product_id"}...)))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_dataset" INNER JOIN "dp_product_dataset"`)).
			WillReturnRows(sqlmock.NewRows(append(dataset.DatasetCols, []string{"dataset_id", "product_id"}...)))

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("get catalog list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(cataloglist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(CatalogCols).
				AddRow(1, time.Now(), time.Now(), nil, cataloglist[0]["title"], cataloglist[0]["description"], cataloglist[0]["featured_medium_id"], cataloglist[0]["published_date"]).
				AddRow(2, time.Now(), time.Now(), nil, cataloglist[1]["title"], cataloglist[1]["description"], cataloglist[1]["featured_medium_id"], cataloglist[1]["published_date"]))

		medium.MediumSelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product" INNER JOIN "dp_catalog_product"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows(append(product.ProductCols, []string{"product_id", "catalog_id"}...)).
				AddRow(1, time.Now(), time.Now(), nil, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1, 1))

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		delete(cataloglist[0], "product_ids")

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(cataloglist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(cataloglist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("get catalog list with paiganation", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(cataloglist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(CatalogCols).
				AddRow(2, time.Now(), time.Now(), nil, cataloglist[1]["title"], cataloglist[1]["description"], cataloglist[1]["featured_medium_id"], cataloglist[1]["published_date"]))

		medium.MediumSelectMock(mock)

		productsAssociationSelectMock(mock, 2)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		delete(cataloglist[1], "product_ids")

		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(cataloglist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(cataloglist[1])

		test.ExpectationsMet(t, mock)
	})
}
