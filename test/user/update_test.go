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

func TestUpdateUser(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	updatedUser := map[string]interface{}{
		"email":      "updatedUser@mail.com",
		"first_name": "Updated User Fname",
		"last_name":  "Updated User LName",
	}

	selectQuery := regexp.QuoteMeta(`SELECT * FROM "dp_user"`)

	t.Run("update user", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "first_name", "last_name"}).
				AddRow(1, time.Now(), time.Now(), nil, "user@mail.com", "User Fname", "User Lname"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_user\" SET (.+)  WHERE (.+) \"dp_user\".\"id\" = `).
			WithArgs(updatedUser["email"], updatedUser["first_name"], updatedUser["last_name"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "first_name", "last_name"}).
				AddRow(1, time.Now(), time.Now(), nil, updatedUser["email"], updatedUser["first_name"], updatedUser["last_name"]))

		e.PUT("/users/1").
			WithJSON(updatedUser).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(updatedUser)

		mock.ExpectationsWereMet()
	})

	t.Run("user record not found", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "first_name", "last_name"}))

		e.PUT("/users/1").
			WithJSON(updatedUser).
			Expect().
			Status(http.StatusNotFound)

		mock.ExpectationsWereMet()
	})

	t.Run("invalid user id", func(t *testing.T) {
		e.PUT("/users/abc").
			WithJSON(updatedUser).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("unprocessable user body", func(t *testing.T) {

		invalidBody := map[string]interface{}{
			"emai":      "updatedUser@mail.com",
			"firt_name": "Updated User Fname",
			"lastname":  "Updated User LName",
		}

		e.PUT("/users/1").
			WithJSON(invalidBody).
			Expect().
			Status(http.StatusUnprocessableEntity)

	})
}
