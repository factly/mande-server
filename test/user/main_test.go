package user

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var user map[string]interface{} = map[string]interface{}{
	"email":      "user@mail.com",
	"first_name": "User Fname",
	"last_name":  "User LName",
}

var invalidUser map[string]interface{} = map[string]interface{}{
	"emil":      "user@mail.com",
	"firs_name": "User Fname",
	"lst_name":  "User LName",
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

var userCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "email", "first_name", "last_name"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_user"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_user"`)

var basePath string = "/users"
var path string = "/users/{user_id}"

func userSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(userCols).
			AddRow(1, time.Now(), time.Now(), nil, user["email"], user["first_name"], user["last_name"]))
}

func userCartExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_cart`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
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
