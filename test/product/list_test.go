package product

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
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
)

func TestListProduct(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	CommonListTests(t, mock, adminExpect)

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	CommonListTests(t, mock, userExpect)

	server.Close()
}

func CommonListTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get empty product list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(ProductCols))

		EmptyProductAssociationsMock(mock)

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("get product list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(productlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(ProductCols).
				AddRow(1, time.Now(), time.Now(), nil, productlist[0]["title"], productlist[0]["slug"], productlist[0]["price"], productlist[0]["status"], productlist[0]["currency_id"], productlist[0]["featured_medium_id"]).
				AddRow(2, time.Now(), time.Now(), nil, productlist[1]["title"], productlist[1]["slug"], productlist[1]["price"], productlist[1]["status"], productlist[1]["currency_id"], productlist[1]["featured_medium_id"]))

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_product_tag"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows(append(tag.TagCols, []string{"tag_id", "product_id"}...)).
				AddRow(1, time.Now(), time.Now(), nil, tag.Tag["title"], tag.Tag["slug"], 1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_dataset" INNER JOIN "dp_product_dataset"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows(append(dataset.DatasetCols, []string{"dataset_id", "product_id"}...)).
				AddRow(1, time.Now(), time.Now(), nil, dataset.Dataset["title"], dataset.Dataset["description"], dataset.Dataset["source"], dataset.Dataset["frequency"], dataset.Dataset["temporal_coverage"], dataset.Dataset["granularity"], dataset.Dataset["contact_name"], dataset.Dataset["contact_email"], dataset.Dataset["license"], dataset.Dataset["data_standard"], dataset.Dataset["sample_url"], dataset.Dataset["related_articles"], dataset.Dataset["time_saved"], dataset.Dataset["price"], dataset.Dataset["currency_id"], dataset.Dataset["featured_medium_id"], 1, 1))

		delete(productlist[0], "tag_ids")
		delete(productlist[0], "dataset_ids")
		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(productlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(productlist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("get product list with paiganation", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(productlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(ProductCols).
				AddRow(2, time.Now(), time.Now(), nil, productlist[1]["title"], productlist[1]["slug"], productlist[1]["price"], productlist[1]["status"], productlist[1]["currency_id"], productlist[1]["featured_medium_id"]))

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tagsAssociationSelectMock(mock, 2)

		datasetsAssociationSelectMock(mock, 2)

		delete(productlist[1], "tag_ids")
		delete(productlist[1], "dataset_ids")
		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(productlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(productlist[1])

		test.ExpectationsMet(t, mock)
	})
}
