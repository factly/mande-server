package currency

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
)

func TestCreateCurrency(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("create a currency", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_currency"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, currency["iso_code"], currency["name"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
		mock.ExpectCommit()

		currencySelectMock(mock)

		e.POST(basePath).
			WithJSON(currency).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(currency)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable currency body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidCurrency).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty currency body", func(t *testing.T) {
		e.POST(basePath).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("create more than one currency", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))

		e.POST(basePath).
			WithJSON(currency).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})
}
