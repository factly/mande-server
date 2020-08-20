package user

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

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

	t.Run("delete user", func(t *testing.T) {
		userSelectMock(mock)

		userCartExpect(mock, 0)

		userMembershipExpect(mock, 0)

		userOrderExpect(mock, 0)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_user" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(path).
			WithPath("user_id", "1").
			Expect().
			Status(http.StatusOK)

		test.ExpectationsMet(t, mock)
	})

	t.Run("user record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(userCols))

		e.DELETE(path).
			WithPath("user_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid user id", func(t *testing.T) {
		e.DELETE(path).
			WithPath("user_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("user associated with cart", func(t *testing.T) {
		userSelectMock(mock)

		userCartExpect(mock, 1)

		e.DELETE(path).
			WithPath("user_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("user associated with membership", func(t *testing.T) {
		userSelectMock(mock)

		userCartExpect(mock, 0)

		userMembershipExpect(mock, 1)

		e.DELETE(path).
			WithPath("user_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("user associated with order", func(t *testing.T) {
		userSelectMock(mock)

		userCartExpect(mock, 0)

		userMembershipExpect(mock, 0)

		userOrderExpect(mock, 1)

		e.DELETE(path).
			WithPath("user_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

}
