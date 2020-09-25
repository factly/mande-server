package order

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/cart"
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/factly/data-portal-server/test/product"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
)

var Order map[string]interface{} = map[string]interface{}{
	"user_id":           1,
	"status":            "teststatus",
	"payment_id":        1,
	"razorpay_order_id": "razor_001",
}

var orderlist []map[string]interface{} = []map[string]interface{}{
	{
		"user_id":           1,
		"status":            "teststatus1",
		"payment_id":        1,
		"razorpay_order_id": "razor_001",
	},
	{
		"user_id":           1,
		"status":            "teststatus2",
		"payment_id":        1,
		"razorpay_order_id": "razor_002",
	},
}

var OrderCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "user_id", "status", "payment_id", "razorpay_order_id"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_order"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_order"`)

const basePath string = "/orders"
const path string = "/orders/{order_id}"

func OrderSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(OrderCols).
			AddRow(1, time.Now(), time.Now(), nil, Order["user_id"], Order["status"], Order["payment_id"], Order["razorpay_order_id"]))
}

func selectWithTwoArgs(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows(OrderCols).
			AddRow(1, time.Now(), time.Now(), nil, Order["user_id"], Order["status"], Order["payment_id"], Order["razorpay_order_id"]))
}

func associatedProductsSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product" INNER JOIN "dp_order_item"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows(append(product.ProductCols, []string{"product_id", "order_id"}...)).
			AddRow(1, time.Now(), time.Now(), nil, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1, 1))
}

func insertMock(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()

	cart.CartItemSelectMock(mock)

	product.ProductSelectMock(mock)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_cart_item" SET "deleted_at"=`)).
		WithArgs(test.AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(`INSERT INTO "dp_order"`).
		WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, "created").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(`SELECT "payment_id", "razorpay_order_id"`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"payment_id", "razorpay_order_id"}).AddRow(nil, nil))

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "dp_order_item"`)).
		WithArgs(1, 1, 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	OrderSelectMock(mock)

	associatedProductsSelectMock(mock)

	dataset.DatasetSelectMock(mock)

	tag.TagSelectMock(mock)
}

func validateAssociations(result *httpexpect.Object) {
	result.Value("payment").
		Object().
		ContainsMap(payment.PaymentReceive)
}
