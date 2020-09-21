package user

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var User map[string]interface{} = map[string]interface{}{
	"email":      "user@mail.com",
	"first_name": "User Fname",
	"last_name":  "User LName",
}

var invalidUser map[string]interface{} = map[string]interface{}{
	"emil":      "user@mail.com",
	"firs_name": "User Fname",
	"lst_name":  "User LName",
}

var undecodableUser map[string]interface{} = map[string]interface{}{
	"email":      5,
	"first_name": 13,
	"last_name":  "User LName",
}

var userlist []map[string]interface{} = []map[string]interface{}{
	{
		"email":      "user1@mail.com",
		"first_name": "User1 FName",
		"last_name":  "User1 LName",
	}, {
		"email":      "user2@mail.com",
		"first_name": "User2 FName",
		"last_name":  "User2 LName",
	},
}

var UserCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "email", "first_name", "last_name"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_user"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_user"`)

const basePath string = "/users"
const path string = "/users/{user_id}"

func UserSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(UserCols).
			AddRow(1, time.Now(), time.Now(), nil, User["email"], User["first_name"], User["last_name"]))
}

func userMembershipExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_membership`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func userOrderExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_order`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}
