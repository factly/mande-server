package currency

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateCurrency(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()

	t.Run("create a currency", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_currency"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Currency["iso_code"], Currency["name"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
		mock.ExpectCommit()

		e.POST(basePath).
			WithJSON(Currency).
			WithHeaders(headers).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(Currency)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable currency body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidCurrency).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty currency body", func(t *testing.T) {
		e.POST(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("createing currency fails", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_currency"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Currency["iso_code"], Currency["name"]).
			WillReturnError(errors.New("cannot create currency"))
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Currency).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create more than one currency", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))

		e.POST(basePath).
			WithJSON(Currency).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})
}
