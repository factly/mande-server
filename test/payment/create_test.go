package payment

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestCreatePayment(t *testing.T) {
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

	t.Run("create a payment", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_payment"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Payment["amount"], Payment["gateway"], Payment["currency_id"], Payment["status"], Payment["razorpay_payment_id"], Payment["razorpay_signature"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		PaymentSelectMock(mock)
		currency.CurrencySelectMock(mock)
		mock.ExpectCommit()

		e.POST(basePath).
			WithJSON(Payment).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(Payment)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable payment body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidPayment).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty payment body", func(t *testing.T) {
		e.POST(basePath).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("currency does not exist", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_payment"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Payment["amount"], Payment["gateway"], Payment["currency_id"], Payment["status"], Payment["razorpay_payment_id"], Payment["razorpay_signature"]).
			WillReturnError(errPaymentCurrencyFK)
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Payment).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a payment when meili is down", func(t *testing.T) {
		gock.Off()
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_payment"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Payment["amount"], Payment["gateway"], Payment["currency_id"], Payment["status"], Payment["razorpay_payment_id"], Payment["razorpay_signature"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		PaymentSelectMock(mock)
		currency.CurrencySelectMock(mock)
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Payment).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
