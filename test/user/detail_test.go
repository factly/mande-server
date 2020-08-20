package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

	t.Run("get user by id", func(t *testing.T) {
		userSelectMock(mock)

		e.GET(path).
			WithPath("user_id", "1").
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

		e.GET(path).
			WithPath("user_id", "1").
			Expect().
			Status(http.StatusNotFound)

		mock.ExpectationsWereMet()

	})

	t.Run("invalid user id", func(t *testing.T) {
		e.GET(path).
			WithPath("user_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})

}
