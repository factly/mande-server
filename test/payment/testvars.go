package payment

import (
	"database/sql/driver"
	"errors"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var headers = map[string]string{
	"X-User":         "1",
	"X-Organisation": "1",
}

// order_FjYVOJ8Vod4lmT

var PaymentOrderReq map[string]interface{} = map[string]interface{}{
	"gateway":             "testgateway.com",
	"currency_id":         1,
	"status":              "complete",
	"for":                 "order",
	"entity_id":           1,
	"razorpay_payment_id": "pay_FjYWQFwuiE89Xp",
	"razorpay_signature":  "114d1d336c37685092345249943acb3007e4c5f8a4fc6493c02c676c7c4e2a0c",
}

var PaymentMembershipReq map[string]interface{} = map[string]interface{}{
	"gateway":             "testgateway.com",
	"currency_id":         1,
	"status":              "complete",
	"for":                 "membership",
	"entity_id":           1,
	"razorpay_payment_id": "pay_FjYWQFwuiE89Xp",
	"razorpay_signature":  "114d1d336c37685092345249943acb3007e4c5f8a4fc6493c02c676c7c4e2a0c",
}

var InvalidSigPaymentReq map[string]interface{} = map[string]interface{}{
	"gateway":             "testgateway.com",
	"currency_id":         1,
	"status":              "complete",
	"for":                 "order",
	"entity_id":           1,
	"razorpay_payment_id": "pay_FjYWQFwuiE89Xp",
	"razorpay_signature":  "114d1d336c37685092345249943acb3007e4c5f8a4fc6493c02c676c7c4e2a",
}

var InvalidForPaymentReq map[string]interface{} = map[string]interface{}{
	"gateway":             "testgateway.com",
	"currency_id":         1,
	"status":              "complete",
	"for":                 "orr",
	"entity_id":           1,
	"razorpay_payment_id": "pay_FjYWQFwuiE89Xp",
	"razorpay_signature":  "114d1d336c37685092345249943acb3007e4c5f8a4fc6493c02c676c7c4e2a0c",
}
var Payment map[string]interface{} = map[string]interface{}{
	"amount":              100,
	"gateway":             "testgateway.com",
	"currency_id":         1,
	"status":              "complete",
	"razorpay_payment_id": "pay_FjYWQFwuiE89Xp",
	"razorpay_signature":  "114d1d336c37685092345249943acb3007e4c5f8a4fc6493c02c676c7c4e2a0c",
}

var undecodablePayment map[string]interface{} = map[string]interface{}{
	"amount":      "100",
	"gateway":     500,
	"currency_id": "1",
	"status":      20,
}

var invalidPayment map[string]interface{} = map[string]interface{}{
	"amt":         100,
	"gateway":     "testgateway.com",
	"currency_id": 0,
	"status":      "sucessful",
}

var paymentlist []map[string]interface{} = []map[string]interface{}{
	{
		"amount":              100,
		"gateway":             "testgateway1.com",
		"currency_id":         1,
		"status":              "sucessful1",
		"razorpay_order_id":   "order001",
		"razorpay_payment_id": "payment001",
		"razorpay_signature":  "signature001",
	},
	{
		"amount":              200,
		"gateway":             "testgateway2.com",
		"currency_id":         1,
		"status":              "sucessful2",
		"razorpay_order_id":   "order002",
		"razorpay_payment_id": "payment002",
		"razorpay_signature":  "signature002",
	},
}

var PaymentCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "amount", "gateway", "currency_id", "status", "razorpay_payment_id", "razorpay_signature"}
var orderCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "user_id", "status", "payment_id", "razorpay_order_id"}
var membershipCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "status", "user_id", "payment_id", "plan_id", "razorpay_order_id"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_payment"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(1) FROM "dp_payment"`)
var errPaymentCurrencyFK = errors.New("pq: insert or update on table \"dp_payment\" violates foreign key constraint \"dp_payment_currency_id_dp_currency_id_foreign\"")

const basePath string = "/payments"
const path string = "/payments/{payment_id}"

func PaymentSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(selectQuery).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(PaymentCols).
			AddRow(1, time.Now(), time.Now(), nil, Payment["amount"], Payment["gateway"], Payment["currency_id"], Payment["status"], Payment["razorpay_payment_id"], Payment["razorpay_signature"]))
}

func paymentOrderExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "dp_order"  WHERE`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func paymentMembershipExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "dp_membership"  WHERE`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}
