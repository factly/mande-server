package order

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/cart"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateOrder(t *testing.T) {
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

	t.Run("create a order", func(t *testing.T) {
		insertMock(mock, nil)

		OrderSelectMock(mock)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		cart.CartSelectMock(mock)

		mock.ExpectCommit()

		result := e.POST(basePath).
			WithJSON(Order).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(Order)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable order body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidOrder).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty order body", func(t *testing.T) {
		e.POST(basePath).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("payment does not exist", func(t *testing.T) {
		insertMock(mock, errOrderPaymentFK)
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Order).
			Expect().
			Status(http.StatusInternalServerError)
	})

	t.Run("user does not exist", func(t *testing.T) {
		insertMock(mock, errOrderUserFK)
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Order).
			Expect().
			Status(http.StatusInternalServerError)
	})

	t.Run("cart does not exist", func(t *testing.T) {
		insertMock(mock, errOrderCartFK)
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Order).
			Expect().
			Status(http.StatusInternalServerError)
	})

	t.Run("create a order when meili is down", func(t *testing.T) {
		gock.Off()
		insertMock(mock, nil)

		OrderSelectMock(mock)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		cart.CartSelectMock(mock)

		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Order).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
