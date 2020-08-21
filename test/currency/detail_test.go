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

func TestDetailCurrency(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get currency by id", func(t *testing.T) {
		CurrencySelectMock(mock)

		e.GET(path).
			WithPath("currency_id", "1").
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

		e.GET(path).
			WithPath("currency_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("currency id invalid", func(t *testing.T) {
		e.GET(path).
			WithPath("currency_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})
}
