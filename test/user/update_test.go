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
	"gopkg.in/h2non/gock.v1"
)

func TestUpdateUser(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("update user", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(UserCols).
				AddRow(1, time.Now(), time.Now(), nil, "user@mail.com", "User Fname", "User Lname"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_user\" SET (.+)  WHERE (.+) \"dp_user\".\"id\" = `).
			WithArgs(User["email"], User["first_name"], User["last_name"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		UserSelectMock(mock)
		mock.ExpectCommit()

		e.PUT(path).
			WithPath("user_id", "1").
			WithJSON(User).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(User)

		test.ExpectationsMet(t, mock)
	})

	t.Run("user record not found", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(UserCols))

		e.PUT(path).
			WithPath("user_id", "1").
			WithJSON(User).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid user id", func(t *testing.T) {
		e.PUT(path).
			WithPath("user_id", "abc").
			WithJSON(User).
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

	t.Run("update user", func(t *testing.T) {
		gock.Off()
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(UserCols).
				AddRow(1, time.Now(), time.Now(), nil, "user@mail.com", "User Fname", "User Lname"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_user\" SET (.+)  WHERE (.+) \"dp_user\".\"id\" = `).
			WithArgs(User["email"], User["first_name"], User["last_name"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		UserSelectMock(mock)
		mock.ExpectRollback()

		e.PUT(path).
			WithPath("user_id", "1").
			WithJSON(User).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
