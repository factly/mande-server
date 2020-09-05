package orderitem

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
)

func TestDeleteOrderItem(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("delete order item", func(t *testing.T) {
		OrderItemSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_order_item" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(path).
			WithPathObject(map[string]interface{}{
				"order_id": "1",
				"item_id":  "1",
			}).
			Expect().
			Status(http.StatusOK)

		test.ExpectationsMet(t, mock)
	})

	t.Run("order item record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderItemCols))

		e.DELETE(path).
			WithPathObject(map[string]interface{}{
				"order_id": "1",
				"item_id":  "1",
			}).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid order item id", func(t *testing.T) {
		e.DELETE(path).
			WithPathObject(map[string]interface{}{
				"order_id": "1",
				"item_id":  "abc",
			}).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("invalid order id", func(t *testing.T) {
		e.DELETE(path).
			WithPathObject(map[string]interface{}{
				"order_id": "abc",
				"item_id":  "1",
			}).
			Expect().
			Status(http.StatusNotFound)
	})
}
