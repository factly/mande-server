package payment

import (
	"errors"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var Payment map[string]interface{} = map[string]interface{}{
	"amount":      100,
	"gateway":     "testgateway.com",
	"currency_id": 1,
	"status":      "sucessful",
}

var invalidPayment map[string]interface{} = map[string]interface{}{
	"amt":         100,
	"gateway":     "testgateway.com",
	"currency_id": 0,
	"status":      "sucessful",
}

var paymentlist []map[string]interface{} = []map[string]interface{}{
	{
		"amount":      100,
		"gateway":     "testgateway1.com",
		"currency_id": 1,
		"status":      "sucessful1",
	},
	{
		"amount":      200,
		"gateway":     "testgateway2.com",
		"currency_id": 1,
		"status":      "sucessful2",
	},
}

var PaymentCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "amount", "gateway", "currency_id", "status"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_payment"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_payment"`)
var errPaymentCurrencyFK = errors.New("pq: insert or update on table \"dp_payment\" violates foreign key constraint \"dp_payment_currency_id_dp_currency_id_foreign\"")

const basePath string = "/payments"
const path string = "/payments/{payment_id}"

func PaymentSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(PaymentCols).
			AddRow(1, time.Now(), time.Now(), nil, Payment["amount"], Payment["gateway"], Payment["currency_id"], Payment["status"]))
}

func paymentOrderExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_order"  WHERE`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func paymentMembershipExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_membership"  WHERE`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}
