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

func TestListCart(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get empty cart list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(CartCols))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product" INNER JOIN "dp_cart_item"`)).
			WillReturnRows(sqlmock.NewRows(append(product.ProductCols, []string{"product_id", "cart_id"}...)))

		product.EmptyProductAssociationsMock(mock)

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("get cart list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(cartlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(CartCols).
				AddRow(1, time.Now(), time.Now(), nil, cartlist[0]["status"], cartlist[0]["user_id"]).
				AddRow(2, time.Now(), time.Now(), nil, cartlist[1]["status"], cartlist[1]["user_id"]))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product" INNER JOIN "dp_cart_item"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows(append(product.ProductCols, []string{"product_id", "cart_id"}...)).
				AddRow(1, time.Now(), time.Now(), nil, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1, 1))

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		delete(cartlist[0], "product_ids")

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(cartlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(cartlist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("get cart list with paiganation", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(cartlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(CartCols).
				AddRow(2, time.Now(), time.Now(), nil, cartlist[1]["status"], cartlist[1]["user_id"]))

		productsAssociationSelectMock(mock, 2)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		delete(cartlist[1], "product_ids")

		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(cartlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(cartlist[1])

		test.ExpectationsMet(t, mock)

	})
}
