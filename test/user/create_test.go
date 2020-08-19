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

func TestCreateUser(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	createdUser := map[string]interface{}{
		"email":      "user@mail.com",
		"first_name": "User Fname",
		"last_name":  "User LName",
	}

	t.Run("create a user", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_user"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, createdUser["email"], createdUser["first_name"], createdUser["last_name"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_user"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "first_name", "last_name"}).
				AddRow(1, time.Now(), time.Now(), nil, createdUser["email"], createdUser["first_name"], createdUser["last_name"]))

		e.POST("/users").
			WithJSON(createdUser).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(createdUser)

		mock.ExpectationsWereMet()

	})

	t.Run("unprocessable user body", func(t *testing.T) {
		invalidUser := map[string]interface{}{
			"emil":      "user@mail.com",
			"firs_name": "User Fname",
			"lst_name":  "User LName",
		}

		e.POST("/users").
			WithJSON(invalidUser).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty user body", func(t *testing.T) {
		e.POST("/users").
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

}
