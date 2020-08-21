package tag

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/test"

	"github.com/factly/data-portal-server/action"
	"github.com/gavv/httpexpect"
)

func TestListTag(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get empty list of tags", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)).
			WillReturnRows(sqlmock.NewRows(TagCols))

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("get non-empty list of tags", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(taglist)))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)).
			WillReturnRows(sqlmock.NewRows(TagCols).
				AddRow(1, time.Now(), time.Now(), nil, taglist[0]["title"], taglist[0]["slug"]).
				AddRow(2, time.Now(), time.Now(), nil, taglist[1]["title"], taglist[1]["slug"]))

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(taglist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(taglist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("get list of tags with paiganation", func(t *testing.T) {

		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(taglist)))

		mock.ExpectQuery(`SELECT \* FROM "dp_tag" (.+) LIMIT 1 OFFSET 1`).
			WillReturnRows(sqlmock.NewRows(TagCols).
				AddRow(2, time.Now(), time.Now(), nil, taglist[1]["title"], taglist[1]["slug"]))

		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(taglist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(taglist[1])

		test.ExpectationsMet(t, mock)
	})

}
