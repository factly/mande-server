package user

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
)

func TestDeleteUser(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	deletedUser := map[string]interface{}{
		"email":      "user@mail.com",
		"first_name": "User Fname",
		"last_name":  "User LName",
	}

	userCols := []string{"id", "created_at", "updated_at", "deleted_at", "email", "first_name", "last_name"}
	selectQuery := regexp.QuoteMeta(`SELECT * FROM "dp_user"`)
	userCartQuery := regexp.QuoteMeta(`SELECT count(*) FROM "dp_cart`)
	userMembershipQuery := regexp.QuoteMeta(`SELECT count(*) FROM "dp_membership`)
	userOrderQuery := regexp.QuoteMeta(`SELECT count(*) FROM "dp_order`)

	t.Run("delete user", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(userCols).
				AddRow(1, time.Now(), time.Now(), nil, deletedUser["email"], deletedUser["first_name"], deletedUser["last_name"]))

		mock.ExpectQuery(userCartQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(userMembershipQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(userOrderQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_user" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE("/users/1").
			Expect().
			Status(http.StatusOK)

		mock.ExpectationsWereMet()

	})

	t.Run("user record not found", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(userCols))

		e.DELETE("/users/1").
			Expect().
			Status(http.StatusNotFound)

		mock.ExpectationsWereMet()
	})

	t.Run("invalid user id", func(t *testing.T) {
		e.DELETE("/users/abc").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("user associated with cart", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(userCols).
				AddRow(1, time.Now(), time.Now(), nil, deletedUser["email"], deletedUser["first_name"], deletedUser["last_name"]))

		mock.ExpectQuery(userCartQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))

		e.DELETE("/users/1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		mock.ExpectationsWereMet()
	})

	t.Run("user associated with membership", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(userCols).
				AddRow(1, time.Now(), time.Now(), nil, deletedUser["email"], deletedUser["first_name"], deletedUser["last_name"]))

		mock.ExpectQuery(userCartQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(userMembershipQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))

		e.DELETE("/users/1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		mock.ExpectationsWereMet()
	})

	t.Run("user associated with order", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(userCols).
				AddRow(1, time.Now(), time.Now(), nil, deletedUser["email"], deletedUser["first_name"], deletedUser["last_name"]))

		mock.ExpectQuery(userCartQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(userMembershipQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(userOrderQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))

		e.DELETE("/users/1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		mock.ExpectationsWereMet()
	})

}
