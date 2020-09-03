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
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

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
		CurrencySelectMock(mock)
		mock.ExpectCommit()

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

	t.Run("update currency when meili is down", func(t *testing.T) {
		gock.Off()
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CurrencyCols).
				AddRow(1, time.Now(), time.Now(), nil, "iso_code", "name"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_currency\" SET (.+)  WHERE (.+) \"dp_currency\".\"id\" = `).
			WithArgs(Currency["iso_code"], Currency["name"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		CurrencySelectMock(mock)
		mock.ExpectRollback()

		e.PUT(path).
			WithPath("currency_id", "1").
			WithJSON(Currency).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
