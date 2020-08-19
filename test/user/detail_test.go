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

func TestDetailUser(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	selectQuery := regexp.QuoteMeta(`SELECT * FROM "dp_user"`)
	user := map[string]interface{}{
		"email":      "user@mail.com",
		"first_name": "User Fname",
		"last_name":  "User LName",
	}

	t.Run("get user by id", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "first_name", "last_name"}).
				AddRow(1, time.Now(), time.Now(), nil, user["email"], user["first_name"], user["last_name"]))

		e.GET("/users/1").
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
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "first_name", "last_name"}))

		e.GET("/users/1").
			Expect().
			Status(http.StatusNotFound)

		mock.ExpectationsWereMet()

	})

	t.Run("invalid user id", func(t *testing.T) {
		e.GET("/users/abc").
			Expect().
			Status(http.StatusNotFound)
	})

}
