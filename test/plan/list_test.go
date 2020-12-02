package plan

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestListPlan(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	test.KavachGock()
	test.KetoGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	CommonListTest(t, mock, adminExpect)

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	gock.New(server.URL).EnableNetworking().Persist()

	CommonListTest(t, mock, userExpect)

	server.Close()
}

func CommonListTest(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get empty plan list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(PlanCols))

		e.GET(basePath).
			WithHeaders(headers).
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
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, planlist[0]["name"], planlist[0]["description"], planlist[0]["status"], planlist[0]["duration"], planlist[0]["price"], planlist[0]["currency_id"]).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, planlist[1]["name"], planlist[1]["description"], planlist[1]["status"], planlist[1]["duration"], planlist[1]["price"], planlist[1]["currency_id"]))

		associatedCatalogSelectMock(mock, 1, 2)
		productCatalogAssociationMock(mock)
		currency.CurrencySelectMock(mock)
		datasetsAssociationSelectMock(mock)
		tagsAssociationSelectMock(mock)
		currency.CurrencySelectMock(mock)

		delete(planlist[0], "catalog_ids")
		e.GET(basePath).
			WithHeaders(headers).
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
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, planlist[1]["name"], planlist[1]["description"], planlist[1]["status"], planlist[1]["duration"], planlist[1]["price"], planlist[1]["currency_id"]))

		associatedCatalogSelectMock(mock, 2)
		currency.CurrencySelectMock(mock)

		delete(planlist[1], "catalog_ids")
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
			ContainsMap(map[string]interface{}{"total": len(planlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(planlist[1])

		test.ExpectationsMet(t, mock)
	})
}
