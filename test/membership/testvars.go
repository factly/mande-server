package membership

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/catalog"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/factly/data-portal-server/test/plan"
	"github.com/factly/data-portal-server/test/product"
	"github.com/gavv/httpexpect"
)

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

var MembershipCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "status", "user_id", "payment_id", "plan_id", "razorpay_order_id"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_membership"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_membership"`)

const basePath string = "/memberships"
const path string = "/memberships/{membership_id}"

func MembershipSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(MembershipCols).
			AddRow(1, time.Now(), time.Now(), nil, Membership["status"], Membership["user_id"], Membership["payment_id"], Membership["plan_id"], Membership["razorpay_order_id"]))
}

func selectWithTwoArgsMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows(MembershipCols).
			AddRow(1, time.Now(), time.Now(), nil, Membership["status"], Membership["user_id"], Membership["payment_id"], Membership["plan_id"], Membership["razorpay_order_id"]))
}

func associatedPlansCatalogSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_catalog" INNER JOIN "dp_plan_catalog"`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows(append(catalog.CatalogCols, []string{"plan_id", "catalog_id"}...)).
			AddRow(1, time.Now(), time.Now(), nil, catalog.Catalog["title"], catalog.Catalog["description"], catalog.Catalog["featured_medium_id"], catalog.Catalog["published_date"], 1, 1))
}

func productCatalogAssociationMock(mock sqlmock.Sqlmock, catId uint) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product" INNER JOIN "dp_catalog_product"`)).
		WithArgs(catId).
		WillReturnRows(sqlmock.NewRows(append(product.ProductCols, []string{"product_id", "catalog_id"}...)).
			AddRow(1, time.Now(), time.Now(), nil, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1, catId))
}

func createMock(mock sqlmock.Sqlmock) {
	plan.PlanSelectMock(mock)
	mock.ExpectQuery(`INSERT INTO "dp_membership"`).
		WithArgs(test.AnyTime{}, test.AnyTime{}, nil, "created", 1, Membership["plan_id"]).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "payment_id", "razorpay_order_id"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"payment_id", "razorpay_order_id"}).AddRow(nil, nil))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_currency"`)).
		WillReturnRows(sqlmock.NewRows(currency.CurrencyCols).
			AddRow(1, time.Now(), time.Now(), nil, currency.Currency["iso_code"], currency.Currency["name"]))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_membership" SET`)).
		WithArgs(test.AnyTime{}, 1, 1, test.RazorpayOrder["id"], "processing", test.AnyTime{}, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	postCreateSelectMock(mock)
}

func postCreateSelectMock(mock sqlmock.Sqlmock) {
	MembershipSelectMock(mock)

	plan.PlanSelectMock(mock)

	associatedPlansCatalogSelectMock(mock)

	productCatalogAssociationMock(mock, 1)
}

func validateAssociations(result *httpexpect.Object) {
	result.Value("plan").
		Object().
		ContainsMap(plan.PlanReceive)

	result.Value("payment").
		Object().
		ContainsMap(payment.PaymentReceive)
}
