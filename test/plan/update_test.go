package plan

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestUpdatePlan(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("update plan", func(t *testing.T) {
		planUpdateMock(mock)
		mock.ExpectCommit()

		e.PUT(path).
			WithPath("plan_id", "1").
			WithJSON(Plan).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(PlanReceive)

		test.ExpectationsMet(t, mock)
	})

	t.Run("plan record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(PlanCols))

		e.PUT(path).
			WithPath("plan_id", "1").
			WithJSON(Plan).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid plan id", func(t *testing.T) {
		e.PUT(path).
			WithPath("plan_id", "abc").
			WithJSON(Plan).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("unprocessable plan body", func(t *testing.T) {
		e.PUT(path).
			WithPath("plan_id", "1").
			WithJSON(invalidPlan).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("undecodable plan body", func(t *testing.T) {
		e.PUT(path).
			WithPath("plan_id", "1").
			WithJSON(undecodablePlan).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})
	t.Run("update plan when meili is down", func(t *testing.T) {
		gock.Off()
		planUpdateMock(mock)
		mock.ExpectRollback()

		e.PUT(path).
			WithPath("plan_id", "1").
			WithJSON(Plan).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
