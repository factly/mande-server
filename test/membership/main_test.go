package membership

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var membership map[string]interface{} = map[string]interface{}{
	"status":     "Test Status",
	"user_id":    1,
	"payment_id": 1,
	"plan_id":    1,
}

var invalidMembership map[string]interface{} = map[string]interface{}{
	"status":    "Test Status",
	"user_id":   0,
	"paymentid": 0,
	"planid":    1,
}

var membershiplist []map[string]interface{} = []map[string]interface{}{
	{
		"status":     "Test Status 1",
		"user_id":    1,
		"payment_id": 1,
		"plan_id":    1,
	},
	{
		"status":     "Test Status 2",
		"user_id":    2,
		"payment_id": 2,
		"plan_id":    2,
	},
}

var membershipCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "status", "user_id", "payment_id", "plan_id"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_membership"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_membership"`)

const basePath string = "/memberships"
const path string = "/memberships/{membership_id}"

func membershipSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(membershipCols).
			AddRow(1, time.Now(), time.Now(), nil, membership["status"], membership["user_id"], membership["payment_id"], membership["plan_id"]))
}

func userSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "first_name", "last_name"}).
			AddRow(1, time.Now(), time.Now(), nil, "email", "first_name", "last_name"))
}

func paymentSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "amount", "gateway", "currency_id", "status"}).
			AddRow(1, time.Now(), time.Now(), nil, 100, "gateway", 1, "status"))
}

func planSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "plan_info", "plan_name", "status"}).
			AddRow(1, time.Now(), time.Now(), nil, "plan_info", "plan_name", "status"))
}
