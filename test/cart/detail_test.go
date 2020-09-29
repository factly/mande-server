package cart

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/factly/data-portal-server/test/product"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
)

func TestDetailCart(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	// ADMIN tests
	CommonDetailTests(t, mock, adminExpect)

	t.Run("get cart item by id", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CartItemCols).
				AddRow(1, time.Now(), time.Now(), nil, CartItem["status"], CartItem["user_id"], CartItem["product_id"]))

		product.ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		adminExpect.GET(path).
			WithHeader("X-User", "1").
			WithPath("cartitem_id", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(CartItem)

		test.ExpectationsMet(t, mock)
	})

	t.Run("cart item record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CartItemCols))

		adminExpect.GET(path).
			WithHeader("X-User", "1").
			WithPath("cartitem_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	// USER test
	CommonDetailTests(t, mock, userExpect)

	t.Run("get cart item by id", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows(CartItemCols).
				AddRow(1, time.Now(), time.Now(), nil, CartItem["status"], CartItem["user_id"], CartItem["product_id"]))

		product.ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		userExpect.GET(path).
			WithHeader("X-User", "1").
			WithPath("cartitem_id", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(CartItem)

		test.ExpectationsMet(t, mock)
	})

	t.Run("cart item record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows(CartItemCols))

		userExpect.GET(path).
			WithHeader("X-User", "1").
			WithPath("cartitem_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid user header", func(t *testing.T) {
		userExpect.GET(path).
			WithHeader("X-User", "abc").
			WithPath("cartitem_id", "1").
			Expect().
			Status(http.StatusNotFound)
	})

	server.Close()
}

func CommonDetailTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {

	t.Run("invalid cart item id", func(t *testing.T) {
		e.GET(path).
			WithHeader("X-User", "1").
			WithPath("cartitem_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})
}
