package payment

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestDeletePayment(t *testing.T) {

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

	t.Run("delete payment", func(t *testing.T) {
		PaymentSelectMock(mock)

		paymentOrderExpect(mock, 0)

		paymentMembershipExpect(mock, 0)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_payment" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(path).
			WithPath("payment_id", "1").
			Expect().
			Status(http.StatusOK)

		test.ExpectationsMet(t, mock)
	})

	t.Run("payment record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(PaymentCols))

		e.DELETE(path).
			WithPath("payment_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid payment id", func(t *testing.T) {
		e.DELETE(path).
			WithPath("payment_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("payment associated with order", func(t *testing.T) {
		PaymentSelectMock(mock)

		paymentOrderExpect(mock, 1)

		e.DELETE(path).
			WithPath("payment_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("payment associated with membership", func(t *testing.T) {
		PaymentSelectMock(mock)

		paymentOrderExpect(mock, 0)

		paymentMembershipExpect(mock, 1)

		e.DELETE(path).
			WithPath("payment_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("delete payment when meili is down", func(t *testing.T) {
		gock.Off()
		PaymentSelectMock(mock)

		paymentOrderExpect(mock, 0)

		paymentMembershipExpect(mock, 0)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_payment" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectRollback()

		e.DELETE(path).
			WithPath("payment_id", "1").
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
