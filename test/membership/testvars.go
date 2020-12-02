package membership

import (
	"database/sql/driver"
	"regexp"
	"time"

	"github.com/factly/data-portal-server/test/product"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/catalog"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/factly/data-portal-server/test/plan"
	"github.com/gavv/httpexpect"
)

var headers = map[string]string{
	"X-User":         "1",
	"X-Organisation": "1",
}

var Membership map[string]interface{} = map[string]interface{}{
	"status":            "Test Status",
	"user_id":           1,
	"payment_id":        1,
	"plan_id":           1,
	"razorpay_order_id": "",
}

var requestBody map[string]interface{} = map[string]interface{}{
	"plan_id": 1,
}

var invalidMembership map[string]interface{} = map[string]interface{}{
	"plan": 1,
}

var membershiplist []map[string]interface{} = []map[string]interface{}{
	{
		"status":            "Test Status 1",
		"user_id":           1,
		"payment_id":        1,
		"plan_id":           1,
		"razorpay_order_id": "",
	},
	{
		"status":            "Test Status 2",
		"user_id":           1,
		"payment_id":        1,
		"plan_id":           1,
		"razorpay_order_id": "",
	},
}

var MembershipCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "created_by_id", "updated_by_id", "status", "user_id", "payment_id", "plan_id", "razorpay_order_id"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_membership"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(1) FROM "dp_membership"`)

const basePath string = "/memberships"
const path string = "/memberships/{membership_id}"

func MembershipSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(selectQuery).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(MembershipCols).
			AddRow(1, time.Now(), time.Now(), nil, 1, 1, Membership["status"], Membership["user_id"], Membership["payment_id"], Membership["plan_id"], Membership["razorpay_order_id"]))
}

func createMock(mock sqlmock.Sqlmock) {
	plan.PlanSelectMock(mock)
	mock.ExpectQuery(`INSERT INTO "dp_membership"`).
		WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, "created", 1, nil, Membership["plan_id"], "").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_currency"`)).
		WillReturnRows(sqlmock.NewRows(currency.CurrencyCols).
			AddRow(1, time.Now(), time.Now(), nil, 1, 1, currency.Currency["iso_code"], currency.Currency["name"]))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_membership" SET`)).
		WithArgs(1, test.AnyTime{}, test.AnyTime{}, 1, 1, "processing", 1, 1, test.RazorpayOrder["id"], 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	postCreateSelectMock(mock)
}

func postCreateSelectMock(mock sqlmock.Sqlmock) {
	MembershipSelectMock(mock)

	plan.PlanSelectMock(mock)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_plan_catalog"`)).
		WillReturnRows(sqlmock.NewRows([]string{"plan_id", "catalog_id"}).
			AddRow(1, 1))

	catalog.CatalogSelectMock(mock)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_catalog_product"`)).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "catalog_id"}).
			AddRow(1, 1))

	product.ProductSelectMock(mock)
}

func validateAssociations(result *httpexpect.Object) {
	result.Value("plan").
		Object().
		ContainsMap(plan.PlanReceive)

	result.Value("payment").
		Object().
		ContainsMap(payment.Payment)
}
