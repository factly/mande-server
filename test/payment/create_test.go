package payment

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

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
	router := action.RegisterUserRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	test.RazorpayGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("create a order payment", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(orderCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, "status", nil, "order_FjYVOJ8Vod4lmT"))

		mock.ExpectQuery(`INSERT INTO "dp_payment"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Payment["amount"], Payment["gateway"], Payment["currency_id"], Payment["status"], Payment["razorpay_payment_id"], Payment["razorpay_signature"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_order" SET`)).
			WithArgs(1, "complete", test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		PaymentSelectMock(mock)
		currency.CurrencySelectMock(mock)
		mock.ExpectCommit()

		e.POST(basePath).
			WithJSON(PaymentOrderReq).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(Payment)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a membership payment", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_membership"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(membershipCols).
				AddRow(1, time.Now(), time.Now(), nil, "status", 1, nil, 1, "order_FjYVOJ8Vod4lmT"))

		mock.ExpectQuery(`INSERT INTO "dp_payment"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Payment["amount"], Payment["gateway"], Payment["currency_id"], Payment["status"], Payment["razorpay_payment_id"], Payment["razorpay_signature"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_membership" SET`)).
			WithArgs(1, "complete", test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		PaymentSelectMock(mock)
		currency.CurrencySelectMock(mock)
		mock.ExpectCommit()

		e.POST(basePath).
			WithJSON(PaymentMembershipReq).
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

	t.Run("undecodable payment body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(undecodablePayment).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty payment body", func(t *testing.T) {
		e.POST(basePath).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid for field in payment request", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(InvalidForPaymentReq).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("order for payment does not exist", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(orderCols))

		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(PaymentOrderReq).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("membership for payment does not exist", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_membership"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(membershipCols))

		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(PaymentMembershipReq).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("signature validation fails", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(orderCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, "status", nil, "order_FjYVOJ8Vod4lmT"))

		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(InvalidSigPaymentReq).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("currency does not exist", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(orderCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, "status", nil, "order_FjYVOJ8Vod4lmT"))

		mock.ExpectQuery(`INSERT INTO "dp_payment"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Payment["amount"], Payment["gateway"], Payment["currency_id"], Payment["status"], Payment["razorpay_payment_id"], Payment["razorpay_signature"]).
			WillReturnError(errPaymentCurrencyFK)
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(PaymentOrderReq).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("razorpay returns error", func(t *testing.T) {
		gock.Off()
		test.MeiliGock()
		gock.New(server.URL).EnableNetworking().Persist()
		defer gock.DisableNetworking()

		gock.New("https://api.razorpay.com").
			Post("/v1/payments/(.+)").
			Reply(http.StatusInternalServerError)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(orderCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, "status", nil, "order_FjYVOJ8Vod4lmT"))

		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(PaymentOrderReq).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("razorpay returns invalid payment body", func(t *testing.T) {
		gock.Off()
		test.MeiliGock()
		gock.New(server.URL).EnableNetworking().Persist()
		defer gock.DisableNetworking()

		gock.New("https://api.razorpay.com").
			Get("/v1/payments/(.+)").
			Reply(http.StatusOK).
			JSON(map[string]interface{}{
				"currency": "INR",
				"receipt":  "Test Receipt no. 1",
			})

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(orderCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, "status", nil, "order_FjYVOJ8Vod4lmT"))

		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(PaymentOrderReq).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	gock.Off()
	test.MeiliGock()
	test.RazorpayGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	t.Run("creating payment fails", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(orderCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, "status", nil, "order_FjYVOJ8Vod4lmT"))

		mock.ExpectQuery(`INSERT INTO "dp_payment"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Payment["amount"], Payment["gateway"], Payment["currency_id"], Payment["status"], Payment["razorpay_payment_id"], Payment["razorpay_signature"]).
			WillReturnError(errors.New(`cannot create payment`))

		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(PaymentOrderReq).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("updating order fails", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(orderCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, "status", nil, "order_FjYVOJ8Vod4lmT"))

		mock.ExpectQuery(`INSERT INTO "dp_payment"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Payment["amount"], Payment["gateway"], Payment["currency_id"], Payment["status"], Payment["razorpay_payment_id"], Payment["razorpay_signature"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_order" SET`)).
			WithArgs(1, "complete", test.AnyTime{}, 1).
			WillReturnError(errors.New(`cannot update order`))

		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(PaymentOrderReq).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("updating membership fails", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_membership"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(membershipCols).
				AddRow(1, time.Now(), time.Now(), nil, "status", 1, nil, 1, "order_FjYVOJ8Vod4lmT"))

		mock.ExpectQuery(`INSERT INTO "dp_payment"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Payment["amount"], Payment["gateway"], Payment["currency_id"], Payment["status"], Payment["razorpay_payment_id"], Payment["razorpay_signature"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_membership" SET`)).
			WithArgs(1, "complete", test.AnyTime{}, 1).
			WillReturnError(errors.New(`cannot update membership`))

		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(PaymentMembershipReq).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a payment when meili is down", func(t *testing.T) {
		gock.Off()
		test.RazorpayGock()
		gock.New(server.URL).EnableNetworking().Persist()
		defer gock.DisableNetworking()

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(orderCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, "status", nil, "order_FjYVOJ8Vod4lmT"))

		mock.ExpectQuery(`INSERT INTO "dp_payment"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Payment["amount"], Payment["gateway"], Payment["currency_id"], Payment["status"], Payment["razorpay_payment_id"], Payment["razorpay_signature"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_order" SET`)).
			WithArgs(1, "complete", test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		PaymentSelectMock(mock)
		currency.CurrencySelectMock(mock)
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(PaymentOrderReq).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
