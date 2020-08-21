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

func TestCreateUser(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("create a user", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_user"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, User["email"], User["first_name"], User["last_name"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
		mock.ExpectCommit()

		UserSelectMock(mock)

		e.POST(basePath).
			WithJSON(User).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(User)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable user body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidUser).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty user body", func(t *testing.T) {
		e.POST(basePath).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

}
