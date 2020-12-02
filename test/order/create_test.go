package order

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
	"github.com/factly/data-portal-server/test/cart"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/product"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateOrder(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterUserRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	test.RazorpayGock()
	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("create a order", func(t *testing.T) {

		insertMock(mock)
		mock.ExpectCommit()

		e.POST(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(Order)

		test.ExpectationsMet(t, mock)
	})

	t.Run("no cart items found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_cart_item"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(cart.CartItemCols))

		mock.ExpectRollback()

		e.POST(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("deleting cart items fail", func(t *testing.T) {
		mock.ExpectBegin()

		cart.CartItemSelectMock(mock)

		product.ProductSelectMock(mock)

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_cart_item" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnError(errors.New(`cannot delete cart item`))

		mock.ExpectRollback()

		e.POST(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("creating order fails", func(t *testing.T) {
		mock.ExpectBegin()

		cart.CartItemSelectMock(mock)

		product.ProductSelectMock(mock)

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_cart_item" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`INSERT INTO "dp_order"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, 1, "created", nil, "").
			WillReturnError(errors.New(`cannot create order`))

		mock.ExpectRollback()

		e.POST(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("when razorpay gives error", func(t *testing.T) {
		gock.Off()
		test.MeiliGock()
		test.KavachGock()
		test.KetoGock()
		gock.New(server.URL).EnableNetworking().Persist()
		defer gock.DisableNetworking()

		gock.New("https://api.razorpay.com").
			Post("/v1/orders").
			Persist().
			Reply(http.StatusInternalServerError)

		mock.ExpectBegin()

		cart.CartItemSelectMock(mock)

		product.ProductSelectMock(mock)

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_cart_item" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`INSERT INTO "dp_order"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, 1, "created", nil, "").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(`INSERT INTO "dp_product"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1).
			WillReturnRows(sqlmock.NewRows([]string{"fetured_medium_id", "id"}).AddRow(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "dp_order_item"`)).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_currency"`)).
			WillReturnRows(sqlmock.NewRows(currency.CurrencyCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, currency.Currency["iso_code"], currency.Currency["name"]))
		mock.ExpectRollback()

		e.POST(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("when razorpay returns invalid body", func(t *testing.T) {
		gock.Off()
		test.MeiliGock()
		test.KavachGock()
		test.KetoGock()
		gock.New(server.URL).EnableNetworking().Persist()
		defer gock.DisableNetworking()

		gock.New("https://api.razorpay.com").
			Post("/v1/orders").
			Persist().
			Reply(http.StatusOK).
			JSON(map[string]interface{}{
				"currency":      "INR",
				"status":        "captured",
				"order_id":      "order_FjYVOJ8Vod4lmT",
				"invoice_id":    nil,
				"international": false,
				"method":        "card",
			})

		mock.ExpectBegin()

		cart.CartItemSelectMock(mock)

		product.ProductSelectMock(mock)

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_cart_item" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`INSERT INTO "dp_order"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, 1, "created", nil, "").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(`INSERT INTO "dp_product"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1).
			WillReturnRows(sqlmock.NewRows([]string{"fetured_medium_id", "id"}).AddRow(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "dp_order_item"`)).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_currency"`)).
			WillReturnRows(sqlmock.NewRows(currency.CurrencyCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, currency.Currency["iso_code"], currency.Currency["name"]))
		mock.ExpectRollback()

		e.POST(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	gock.Off()
	test.MeiliGock()
	test.RazorpayGock()
	test.KavachGock()
	test.KetoGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	t.Run("updating order fails", func(t *testing.T) {
		mock.ExpectBegin()

		cart.CartItemSelectMock(mock)

		product.ProductSelectMock(mock)

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_cart_item" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`INSERT INTO "dp_order"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, 1, "created", nil, "").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(`INSERT INTO "dp_product"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1).
			WillReturnRows(sqlmock.NewRows([]string{"fetured_medium_id", "id"}).AddRow(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "dp_order_item"`)).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_currency"`)).
			WillReturnRows(sqlmock.NewRows(currency.CurrencyCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, currency.Currency["iso_code"], currency.Currency["name"]))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_order" SET`)).
			WithArgs(1, test.AnyTime{}, test.AnyTime{}, 1, 1, 1, "processing", test.RazorpayOrder["id"], 1).
			WillReturnError(errors.New(`cannot update order`))

		mock.ExpectRollback()

		e.POST(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a order when meili is down", func(t *testing.T) {
		gock.Off()
		test.RazorpayGock()
		test.KavachGock()
		test.KetoGock()
		gock.New(server.URL).EnableNetworking().Persist()
		defer gock.DisableNetworking()

		insertMock(mock)
		mock.ExpectRollback()

		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(Order).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
