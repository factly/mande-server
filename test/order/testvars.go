package order

import (
	"errors"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/cart"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/gavv/httpexpect"
)

var Order map[string]interface{} = map[string]interface{}{
	"user_id":    1,
	"status":     "teststatus",
	"payment_id": 1,
	"cart_id":    1,
}

var invalidOrder map[string]interface{} = map[string]interface{}{
	"userid":    1,
	"status":    "teststatus",
	"paymentid": 1,
	"cartid":    1,
}

var orderlist []map[string]interface{} = []map[string]interface{}{
	{
		"user_id":    1,
		"status":     "teststatus1",
		"payment_id": 1,
		"cart_id":    1,
	},
	{
		"user_id":    1,
		"status":     "teststatus2",
		"payment_id": 1,
		"cart_id":    1,
	},
}

var OrderCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "user_id", "status", "payment_id", "cart_id"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_order"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_order"`)
var errOrderCartFK error = errors.New(`pq: insert or update on table "dp_order" violates foreign key constraint "dp_order_cart_id_dp_cart_id_foreign"`)
var errOrderPaymentFK error = errors.New(`pq: insert or update on table "dp_order" violates foreign key constraint "dp_order_payment_id_dp_payment_id_foreign"`)
var errOrderUserFK error = errors.New(`pq: insert or update on table "dp_order" violates foreign key constraint "dp_order_user_id_dp_user_id_foreign"`)

const basePath string = "/orders"
const path string = "/orders/{order_id}"

func OrderSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(OrderCols).
			AddRow(1, time.Now(), time.Now(), nil, Order["user_id"], Order["status"], Order["payment_id"], Order["cart_id"]))
}

func insertMock(mock sqlmock.Sqlmock, err error) {
	if err != nil {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_order"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Order["user_id"], Order["status"], Order["payment_id"], Order["cart_id"]).
			WillReturnError(err)
	} else {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_order"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Order["user_id"], Order["status"], Order["payment_id"], Order["cart_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	}
}

func updateMock(mock sqlmock.Sqlmock, err error) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(OrderCols).
			AddRow(1, time.Now(), time.Now(), nil, 2, "status", 2, 2))

	mock.ExpectBegin()
	if err != nil {
		mock.ExpectExec(`UPDATE \"dp_order\" SET (.+)  WHERE (.+) \"dp_order\".\"id\" = `).
			WithArgs(Order["cart_id"], Order["payment_id"], Order["status"], test.AnyTime{}, Order["user_id"], 1).
			WillReturnError(err)
	} else {
		mock.ExpectExec(`UPDATE \"dp_order\" SET (.+)  WHERE (.+) \"dp_order\".\"id\" = `).
			WithArgs(Order["cart_id"], Order["payment_id"], Order["status"], test.AnyTime{}, Order["user_id"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func validateAssociations(result *httpexpect.Object) {
	result.Value("cart").
		Object().
		ContainsMap(cart.CartReceive)

	result.Value("payment").
		Object().
		ContainsMap(payment.Payment)
}
