package user

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
)

func TestUpdateUser(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("update user", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(userCols).
				AddRow(1, time.Now(), time.Now(), nil, "user@mail.com", "User Fname", "User Lname"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_user\" SET (.+)  WHERE (.+) \"dp_user\".\"id\" = `).
			WithArgs(user["email"], user["first_name"], user["last_name"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		userSelectMock(mock)

		e.PUT(path).
			WithPath("user_id", "1").
			WithJSON(user).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(user)

		mock.ExpectationsWereMet()
	})

	t.Run("user record not found", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(userCols))

		e.PUT(path).
			WithPath("user_id", "1").
			WithJSON(user).
			Expect().
			Status(http.StatusNotFound)

		mock.ExpectationsWereMet()
	})

	t.Run("invalid user id", func(t *testing.T) {
		e.PUT(path).
			WithPath("user_id", "abc").
			WithJSON(user).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("unprocessable user body", func(t *testing.T) {
		e.PUT(path).
			WithPath("user_id", "1").
			WithJSON(invalidUser).
			Expect().
			Status(http.StatusUnprocessableEntity)

	})
}
