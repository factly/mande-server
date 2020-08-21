package plan

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

func TestUpdatePlan(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("update plan", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(PlanCols).
				AddRow(1, time.Now(), time.Now(), nil, "plan_info", "plan_name", "status"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_plan\" SET (.+)  WHERE (.+) \"dp_plan\".\"id\" = `).
			WithArgs(Plan["plan_info"], Plan["plan_name"], Plan["status"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		PlanSelectMock(mock)

		e.PUT(path).
			WithPath("plan_id", "1").
			WithJSON(Plan).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(Plan)

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
}
