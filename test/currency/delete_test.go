package currency

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestDeleteCurrency(t *testing.T) {

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

	t.Run("delete currency", func(t *testing.T) {
		CurrencySelectMock(mock)

		currencyPaymentExpect(mock, 0)

		currencyProductExpect(mock, 0)

		currencyDatasetExpect(mock, 0)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_currency" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("currency_id", "1").
			Expect().
			Status(http.StatusOK)

		test.ExpectationsMet(t, mock)
	})

	t.Run("currency record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CurrencyCols))

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("currency_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid currency id", func(t *testing.T) {
		e.DELETE(path).
			WithPath("currency_id", "abc").
			WithHeaders(headers).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("currency associated with payment", func(t *testing.T) {
		CurrencySelectMock(mock)

		currencyPaymentExpect(mock, 1)

		e.DELETE(path).
			WithPath("currency_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("currency associated with product", func(t *testing.T) {
		CurrencySelectMock(mock)

		currencyPaymentExpect(mock, 0)

		currencyProductExpect(mock, 1)

		e.DELETE(path).
			WithPath("currency_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("currency associated with dataset", func(t *testing.T) {
		CurrencySelectMock(mock)

		currencyPaymentExpect(mock, 0)

		currencyProductExpect(mock, 0)

		currencyDatasetExpect(mock, 1)

		e.DELETE(path).
			WithPath("currency_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

}
