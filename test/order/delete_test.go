package order

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/product"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestDeleteOrder(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterUserRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("delete order", func(t *testing.T) {
		OrderSelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order_item"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"order_id", "product_id"}).
				AddRow(1, 1))
		product.ProductSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_order_item"`)).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_order" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("order_id", "1").
			Expect().
			Status(http.StatusOK)

		test.ExpectationsMet(t, mock)
	})

	t.Run("order record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderCols))

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("order_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid order id", func(t *testing.T) {
		e.DELETE(path).
			WithHeaders(headers).
			WithPath("order_id", "abc").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("deleting order fail", func(t *testing.T) {
		OrderSelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order_item"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"order_id", "product_id"}).
				AddRow(1, 1))
		product.ProductSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_order_item"`)).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_order" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnError(errors.New("cannot delete order"))
		mock.ExpectRollback()

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("order_id", "1").
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("delete order when meili is down", func(t *testing.T) {
		gock.Off()
		test.KavachGock()
		test.KetoGock()
		gock.New(server.URL).EnableNetworking().Persist()

		OrderSelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order_item"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"order_id", "product_id"}).
				AddRow(1, 1))
		product.ProductSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_order_item"`)).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_order" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectRollback()

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("order_id", "1").
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
