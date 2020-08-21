package plan

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
)

func TestCreatePlan(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("create a plan", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_plan"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Plan["plan_name"], Plan["plan_info"], Plan["status"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
		mock.ExpectCommit()

		PlanSelectMock(mock)

		e.POST(basePath).
			WithJSON(Plan).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(Plan)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable plan body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidPlan).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty plan body", func(t *testing.T) {
		e.POST(basePath).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})
}
