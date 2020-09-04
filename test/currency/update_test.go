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
)

func TestUpdateCurrency(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("update currency", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CurrencyCols).
				AddRow(1, time.Now(), time.Now(), nil, "iso_code", "name"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_currency\" SET (.+)  WHERE (.+) \"dp_currency\".\"id\" = `).
			WithArgs(Currency["iso_code"], Currency["name"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		CurrencySelectMock(mock)

		e.PUT(path).
			WithPath("currency_id", "1").
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
			WithPath("currency_id", "1").
			WithJSON(Currency).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable currency body", func(t *testing.T) {
		e.PUT(path).
			WithPath("currency_id", "1").
			WithJSON(invalidCurrency).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid currency id", func(t *testing.T) {
		e.PUT(path).
			WithPath("currency_id", "abc").
			WithJSON(Currency).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("undecodable currency body", func(t *testing.T) {
		e.PUT(path).
			WithPath("currency_id", "1").
			WithJSON(undecodableCurrency).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})
}
