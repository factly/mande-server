package order

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
)

func TestDetailOrder(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get order by id", func(t *testing.T) {
		OrderSelectMock(mock)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		associatedProductsSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		tag.TagSelectMock(mock)

		result := e.GET(path).
			WithHeader("X-User", "1").
			WithPath("order_id", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(Order)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("order record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderCols))

		e.GET(path).
			WithHeader("X-User", "1").
			WithPath("order_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid order id", func(t *testing.T) {
		e.GET(path).
			WithHeader("X-User", "1").
			WithPath("order_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})
}
