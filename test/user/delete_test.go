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
	"gopkg.in/h2non/gock.v1"
)

func TestDeleteUser(t *testing.T) {

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

	t.Run("delete user", func(t *testing.T) {
		UserSelectMock(mock)

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
			WillReturnRows(sqlmock.NewRows(UserCols))

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
		UserSelectMock(mock)

		userCartExpect(mock, 1)

		e.DELETE(path).
			WithPath("user_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("user associated with membership", func(t *testing.T) {
		UserSelectMock(mock)

		userCartExpect(mock, 0)

		userMembershipExpect(mock, 1)

		e.DELETE(path).
			WithPath("user_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("user associated with order", func(t *testing.T) {
		UserSelectMock(mock)

		userCartExpect(mock, 0)

		userMembershipExpect(mock, 0)

		userOrderExpect(mock, 1)

		e.DELETE(path).
			WithPath("user_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("delete user when meili is down", func(t *testing.T) {
		gock.Off()
		UserSelectMock(mock)

		userCartExpect(mock, 0)

		userMembershipExpect(mock, 0)

		userOrderExpect(mock, 0)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_user" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectRollback()

		e.DELETE(path).
			WithPath("user_id", "1").
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
