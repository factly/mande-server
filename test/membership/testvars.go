package membership

import (
	"errors"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/factly/data-portal-server/test/plan"
	"github.com/factly/data-portal-server/test/user"
	"github.com/gavv/httpexpect"
)

var Membership map[string]interface{} = map[string]interface{}{
	"status":     "Test Status",
	"user_id":    1,
	"payment_id": 1,
	"plan_id":    1,
}

var undecodableMembership map[string]interface{} = map[string]interface{}{
	"status":     10,
	"user_id":    "1",
	"payment_id": 1,
	"plan_id":    "1",
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
		"user_id":    1,
		"payment_id": 1,
		"plan_id":    1,
	},
}

var MembershipCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "status", "user_id", "payment_id", "plan_id"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_membership"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_membership"`)
var errMembershipUserFK error = errors.New(`pq: insert or update on table "dp_membership" violates foreign key constraint "dp_membership_user_id_dp_user_id_foreign"`)
var errMembershipPlanFK error = errors.New(`pq: insert or update on table "dp_membership" violates foreign key constraint "dp_membership_plan_id_dp_plan_id_foreign"`)
var errMembershipPaymentFK error = errors.New(`pq: insert or update on table "dp_membership" violates foreign key constraint "dp_membership_payment_id_dp_payment_id_foreign"`)

const basePath string = "/memberships"
const path string = "/memberships/{membership_id}"

func MembershipSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(MembershipCols).
			AddRow(1, time.Now(), time.Now(), nil, Membership["status"], Membership["user_id"], Membership["payment_id"], Membership["plan_id"]))
}

func updateWithErrorExpect(mock sqlmock.Sqlmock, err error) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(MembershipCols).
			AddRow(1, time.Now(), time.Now(), nil, "status", 2, 2, 2))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE \"dp_membership\" SET (.+)  WHERE (.+) \"dp_membership\".\"id\" = `).
		WithArgs(Membership["payment_id"], Membership["plan_id"], Membership["status"], test.AnyTime{}, Membership["user_id"], 1).
		WillReturnError(err)
	mock.ExpectRollback()
}

func insertWithErrorExpect(mock sqlmock.Sqlmock, err error) {
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "dp_membership"`).
		WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Membership["status"], Membership["user_id"], Membership["payment_id"], Membership["plan_id"]).
		WillReturnError(err)
	mock.ExpectRollback()
}

func validateAssociations(result *httpexpect.Object) {
	result.Value("user").
		Object().
		ContainsMap(user.User)

	result.Value("plan").
		Object().
		ContainsMap(plan.Plan)

	result.Value("payment").
		Object().
		ContainsMap(payment.Payment)
}
