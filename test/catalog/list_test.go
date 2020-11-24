package catalog

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestListCatalog(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	test.MeiliGock()
	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	CommonListTests(t, mock, adminExpect)

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	CommonListTests(t, mock, userExpect)

	server.Close()
}

func CommonListTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get empty catalog list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(CatalogCols))

		e.GET(basePath).
			WithHeaders(headers).
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
		productsAssociationSelectMock(mock, 1, 2)

		currency.CurrencySelectMock(mock)
		datasetsAssociationSelectMock(mock)
		medium.MediumSelectMock(mock)
		tagsAssociationSelectMock(mock)

		delete(cataloglist[0], "product_ids")

		e.GET(basePath).
			WithHeaders(headers).
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

		delete(cataloglist[1], "product_ids")

		e.GET(basePath).
			WithHeaders(headers).
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
