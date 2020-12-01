package cart

import (
	"errors"
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

func TestDeleteCart(t *testing.T) {

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

	t.Run("delete cart item", func(t *testing.T) {
		CartItemSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_cart_item" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("cartitem_id", "1").
			Expect().
			Status(http.StatusOK)

		test.ExpectationsMet(t, mock)
	})

	t.Run("cart record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CartItemCols))

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("cartitem_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid cart id", func(t *testing.T) {
		e.DELETE(path).
			WithHeaders(headers).
			WithPath("cartitem_id", "abc").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("deleting cart items fail", func(t *testing.T) {
		CartItemSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_cart_item" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnError(errors.New("cannot delete cart items"))
		mock.ExpectRollback()

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("cartitem_id", "1").
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("delete cart item when meili is down", func(t *testing.T) {
		gock.Off()
		CartItemSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_cart_item" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectRollback()

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("cartitem_id", "1").
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
