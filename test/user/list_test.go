package user

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

func TestListUser(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get empty users list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(UserCols))

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("list users", func(t *testing.T) {

		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(userlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(UserCols).
				AddRow(1, time.Now(), time.Now(), nil, userlist[0]["email"], userlist[0]["first_name"], userlist[0]["last_name"]).
				AddRow(2, time.Now(), time.Now(), nil, userlist[1]["email"], userlist[1]["first_name"], userlist[1]["last_name"]))

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(userlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(userlist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("list users paiganation", func(t *testing.T) {

		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(userlist)))

		mock.ExpectQuery(`SELECT \* FROM "dp_user" (.+) LIMIT 1 OFFSET 1`).
			WillReturnRows(sqlmock.NewRows(UserCols).
				AddRow(2, time.Now(), time.Now(), nil, userlist[1]["email"], userlist[1]["first_name"], userlist[1]["last_name"]))

		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(userlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(userlist[1])

		test.ExpectationsMet(t, mock)
	})
}
