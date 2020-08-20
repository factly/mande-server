package plan

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var plan map[string]interface{} = map[string]interface{}{
	"plan_name": "Test Plan",
	"plan_info": "Test Plan Info",
	"status":    "teststatus",
}

var invalidPlan map[string]interface{} = map[string]interface{}{
	"planname":  "Test Plan",
	"plan_info": "Test Plan Info",
	"status":    "teststatus",
}

var planlist []map[string]interface{} = []map[string]interface{}{
	{
		"plan_name": "Test Plan 1",
		"plan_info": "Test Plan Info 1",
		"status":    "teststatus1",
	},
	{
		"plan_name": "Test Plan 2",
		"plan_info": "Test Plan Info 2",
		"status":    "teststatus2",
	},
}

var planCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "plan_info", "plan_name", "status"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_plan"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_plan"`)

const basePath string = "/plans"
const path string = "/plans/{plan_id}"

func planSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(planCols).
			AddRow(1, time.Now(), time.Now(), nil, plan["plan_info"], plan["plan_name"], plan["status"]))
}

func planMembershipExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_membership`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}
