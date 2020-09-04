package payment

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

func TestUpdatePayment(t *testing.T) {
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

	t.Run("update payment", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(PaymentCols).
				AddRow(1, time.Now(), time.Now(), nil, 100, "gateway", 1, "status"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_payment\" SET (.+)  WHERE (.+) \"dp_payment\".\"id\" = `).
			WithArgs(Payment["amount"], Payment["currency_id"], Payment["gateway"], Payment["status"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		PaymentSelectMock(mock)
		currency.CurrencySelectMock(mock)
		mock.ExpectCommit()

		e.PUT(path).
			WithPath("payment_id", "1").
			WithJSON(Payment).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(Payment)

		test.ExpectationsMet(t, mock)
	})

	t.Run("payment record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(PaymentCols))

		e.PUT(path).
			WithPath("payment_id", "1").
			WithJSON(Payment).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable payment body", func(t *testing.T) {
		e.PUT(path).
			WithPath("payment_id", "1").
			WithJSON(invalidPayment).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid payment id", func(t *testing.T) {
		e.PUT(path).
			WithPath("payment_id", "abc").
			WithJSON(Payment).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("new currency does not exist", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(PaymentCols).
				AddRow(1, time.Now(), time.Now(), nil, 100, "gateway", 1, "status"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_payment\" SET (.+)  WHERE (.+) \"dp_payment\".\"id\" = `).
			WithArgs(Payment["amount"], Payment["currency_id"], Payment["gateway"], Payment["status"], test.AnyTime{}, 1).
			WillReturnError(errPaymentCurrencyFK)
		mock.ExpectRollback()

		e.PUT(path).
			WithPath("payment_id", "1").
			WithJSON(Payment).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("update payment when meili is down", func(t *testing.T) {
		gock.Off()
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(PaymentCols).
				AddRow(1, time.Now(), time.Now(), nil, 100, "gateway", 1, "status"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_payment\" SET (.+)  WHERE (.+) \"dp_payment\".\"id\" = `).
			WithArgs(Payment["amount"], Payment["currency_id"], Payment["gateway"], Payment["status"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		PaymentSelectMock(mock)
		currency.CurrencySelectMock(mock)
		mock.ExpectRollback()

		e.PUT(path).
			WithPath("payment_id", "1").
			WithJSON(Payment).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})
}
