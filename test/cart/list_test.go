package cart

import (
	"net/http"
	"net/http/httptest"
	"regexp"
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

func TestListCartItems(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	// ADMIN specific tests
	CommonListTests(t, mock, adminExpect)

	t.Run("get cart item list with user query", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(cartitemslist)))

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CartItemCols).
				AddRow(1, time.Now(), time.Now(), nil, cartitemslist[0]["status"], cartitemslist[0]["user_id"], cartitemslist[0]["product_id"]).
				AddRow(2, time.Now(), time.Now(), nil, cartitemslist[1]["status"], cartitemslist[1]["user_id"], cartitemslist[1]["product_id"]))

		product.ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		delete(cartitemslist[0], "product_id")

		adminExpect.GET(basePath).
			WithHeader("X-User", "1").
			WithQuery("user", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(cartitemslist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(cartitemslist[0])

		cartitemslist[0]["product_id"] = 1

		test.ExpectationsMet(t, mock)
	})

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	// USER specific tests
	CommonListTests(t, mock, userExpect)

	t.Run("invalid user header", func(t *testing.T) {
		userExpect.GET(basePath).
			WithHeader("X-User", "anc").
			Expect().
			Status(http.StatusNotFound)
	})

	server.Close()
}

func CommonListTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get empty cart list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(CartItemCols))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_product_tag"`)).
			WillReturnRows(sqlmock.NewRows(tag.TagCols))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_dataset" INNER JOIN "dp_product_dataset"`)).
			WillReturnRows(sqlmock.NewRows(dataset.DatasetCols))

		e.GET(basePath).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("get cart item list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(cartitemslist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(CartItemCols).
				AddRow(1, time.Now(), time.Now(), nil, cartitemslist[0]["status"], cartitemslist[0]["user_id"], cartitemslist[0]["product_id"]).
				AddRow(2, time.Now(), time.Now(), nil, cartitemslist[1]["status"], cartitemslist[1]["user_id"], cartitemslist[1]["product_id"]))

		product.ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		delete(cartitemslist[0], "product_id")

		e.GET(basePath).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(cartitemslist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(cartitemslist[0])

		cartitemslist[0]["product_id"] = 1

		test.ExpectationsMet(t, mock)
	})

	t.Run("get cart item list with paiganation", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(cartitemslist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(CartItemCols).
				AddRow(2, time.Now(), time.Now(), nil, cartitemslist[1]["status"], cartitemslist[1]["user_id"], cartitemslist[1]["product_id"]))

		product.ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		delete(cartitemslist[1], "product_id")

		e.GET(basePath).
			WithHeader("X-User", "1").
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(cartitemslist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(cartitemslist[1])

		cartitemslist[1]["product_id"] = 1

		test.ExpectationsMet(t, mock)

	})
}
