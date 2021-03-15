package order

import (
	"database/sql/driver"
	"regexp"
	"time"

	"github.com/factly/mande-server/test/currency"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/test"
	"github.com/factly/mande-server/test/cart"
	"github.com/factly/mande-server/test/dataset"
	"github.com/factly/mande-server/test/payment"
	"github.com/factly/mande-server/test/product"
	"github.com/factly/mande-server/test/tag"
	"github.com/gavv/httpexpect"
)

var headers = map[string]string{
	"X-User":         "1",
	"X-Organisation": "1",
}

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

var OrderCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "created_by_id", "updated_by_id", "user_id", "status", "payment_id", "razorpay_order_id"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_order"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(1) FROM "dp_order"`)

const basePath string = "/orders"
const path string = "/orders/{order_id}"

func OrderSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(selectQuery).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(OrderCols).
			AddRow(1, time.Now(), time.Now(), nil, 1, 1, Order["user_id"], Order["status"], Order["payment_id"], Order["razorpay_order_id"]))
}

func associatedProductsSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order_item"`)).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows([]string{"order_id", "product_id"}).
			AddRow(1, 1))
	product.ProductSelectMock(mock)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product_dataset"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "dataset_id"}).
			AddRow(1, 1))
	dataset.DatasetSelectMock(mock)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product_tag"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "tag_id"}).
			AddRow(1, 1))
	tag.TagSelectMock(mock)
}

func insertMock(mock sqlmock.Sqlmock) {
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
		WillReturnResult(sqlmock.NewResult(1, 1))

	OrderSelectMock(mock)

	associatedProductsSelectMock(mock)

}

func validateAssociations(result *httpexpect.Object) {
	result.Value("payment").
		Object().
		ContainsMap(payment.Payment)
}
