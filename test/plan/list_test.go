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

func TestListPlan(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get empty plan list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(PlanCols))

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("list plans", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(planlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(PlanCols).
				AddRow(1, time.Now(), time.Now(), nil, planlist[0]["plan_info"], planlist[0]["plan_name"], planlist[0]["status"]).
				AddRow(2, time.Now(), time.Now(), nil, planlist[1]["plan_info"], planlist[1]["plan_name"], planlist[1]["status"]))

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(planlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(planlist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("get plan list with paiganation", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(planlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(PlanCols).
				AddRow(2, time.Now(), time.Now(), nil, planlist[1]["plan_info"], planlist[1]["plan_name"], planlist[1]["status"]))

		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(planlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(planlist[1])

		test.ExpectationsMet(t, mock)
	})
}
