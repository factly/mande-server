package plan

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/catalog"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestCreatePlan(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	test.KavachGock()
	test.KetoGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("create a plan", func(t *testing.T) {
		catalog.CatalogSelectMock(mock)

		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "dp_plan"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Plan["name"], Plan["description"], Plan["price"], Plan["currency_id"], Plan["duration"], Plan["status"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

		mock.ExpectQuery(`INSERT INTO "dp_catalog"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, catalog.Catalog["title"], catalog.Catalog["description"], test.AnyTime{}, catalog.Catalog["featured_medium_id"], 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(`INSERT INTO "dp_plan_catalog"`).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		PlanSelectMock(mock)
		associatedCatalogSelectMock(mock)
		productCatalogAssociationMock(mock, 1)
		currency.CurrencySelectMock(mock)
		datasetsAssociationSelectMock(mock)
		tagsAssociationSelectMock(mock)
		currency.CurrencySelectMock(mock)

		mock.ExpectCommit()

		e.POST(basePath).
			WithJSON(Plan).
			WithHeaders(headers).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(PlanReceive)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable plan body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidPlan).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty plan body", func(t *testing.T) {
		e.POST(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("creating plan fails", func(t *testing.T) {
		catalog.CatalogSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_plan"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Plan["name"], Plan["description"], Plan["price"], Plan["currency_id"], Plan["duration"], Plan["status"]).
			WillReturnError(errors.New("cannot create plan"))
		mock.ExpectRollback()

		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(Plan).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a plan when meili is down", func(t *testing.T) {
		gock.Off()
		test.KavachGock()
		test.KetoGock()
		gock.New(server.URL).EnableNetworking().Persist()

		catalog.CatalogSelectMock(mock)

		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "dp_plan"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Plan["name"], Plan["description"], Plan["price"], Plan["currency_id"], Plan["duration"], Plan["status"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

		mock.ExpectQuery(`INSERT INTO "dp_catalog"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, catalog.Catalog["title"], catalog.Catalog["description"], test.AnyTime{}, catalog.Catalog["featured_medium_id"], 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(`INSERT INTO "dp_plan_catalog"`).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		PlanSelectMock(mock)
		associatedCatalogSelectMock(mock)
		productCatalogAssociationMock(mock, 1)
		currency.CurrencySelectMock(mock)
		datasetsAssociationSelectMock(mock)
		tagsAssociationSelectMock(mock)
		currency.CurrencySelectMock(mock)

		mock.ExpectRollback()

		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(Plan).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
