package currency

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestUpdateCurrency(t *testing.T) {
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

	t.Run("update currency", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CurrencyCols).
				AddRow(1, time.Now(), time.Now(), nil, "iso_code", "name"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_currency\"`).
			WithArgs(test.AnyTime{}, Currency["iso_code"], Currency["name"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		CurrencySelectMock(mock, 1, 1)

		e.PUT(path).
			WithPath("currency_id", "1").
			WithHeaders(headers).
			WithJSON(Currency).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(Currency)

		test.ExpectationsMet(t, mock)
	})

	t.Run("currency record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CurrencyCols))

		e.PUT(path).
			WithHeaders(headers).
			WithPath("currency_id", "1").
			WithJSON(Currency).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable currency body", func(t *testing.T) {
		e.PUT(path).
			WithPath("currency_id", "1").
			WithHeaders(headers).
			WithJSON(invalidCurrency).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid currency id", func(t *testing.T) {
		e.PUT(path).
			WithPath("currency_id", "abc").
			WithHeaders(headers).
			WithJSON(Currency).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("undecodable currency body", func(t *testing.T) {
		e.PUT(path).
			WithPath("currency_id", "1").
			WithHeaders(headers).
			WithJSON(undecodableCurrency).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})
}
