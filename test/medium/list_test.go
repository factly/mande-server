package medium

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
)

func TestListMedium(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	mediumlist := []map[string]interface{}{
		{
			"name":        "Test Medium 1",
			"slug":        "test-medium-1",
			"type":        "testtype1",
			"title":       "Test Title 1",
			"description": "Test Description 1",
			"caption":     "Test Caption 1",
			"alt_text":    "Test alt text 1",
			"file_size":   100,
			"url":         "http:/testurl1.com",
			"dimensions":  "testdims1",
		},
		{
			"name":        "Test Medium 2",
			"slug":        "test-medium-2",
			"type":        "testtype2",
			"title":       "Test Title 2",
			"description": "Test Description 2",
			"caption":     "Test Caption 2",
			"alt_text":    "Test alt text 2",
			"file_size":   200,
			"url":         "http:/testurl2.com",
			"dimensions":  "testdims2",
		},
	}

	mediumCols := []string{"id", "created_at", "updated_at", "deleted_at", "name", "slug", "type", "title", "description", "caption", "alt_text", "file_size", "url", "dimensions"}
	selectQuery := regexp.QuoteMeta(`SELECT * FROM "dp_medium"`)
	countQuery := regexp.QuoteMeta(`SELECT count(*) FROM "dp_medium"`)

	t.Run("empty medium list", func(t *testing.T) {

		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(mediumCols))

		e.GET("/media").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		mock.ExpectationsWereMet()

	})

	t.Run("medium list", func(t *testing.T) {

		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(mediumlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, mediumlist[0]["name"], mediumlist[0]["slug"], mediumlist[0]["type"], mediumlist[0]["title"], mediumlist[0]["description"], mediumlist[0]["caption"], mediumlist[0]["alt_text"], mediumlist[0]["file_size"], mediumlist[0]["url"], mediumlist[0]["dimensions"]).
				AddRow(2, time.Now(), time.Now(), nil, mediumlist[1]["name"], mediumlist[1]["slug"], mediumlist[1]["type"], mediumlist[1]["title"], mediumlist[1]["description"], mediumlist[1]["caption"], mediumlist[1]["alt_text"], mediumlist[1]["file_size"], mediumlist[1]["url"], mediumlist[1]["dimensions"]))

		e.GET("/media").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(mediumlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(mediumlist[0])

		mock.ExpectationsWereMet()

	})

	t.Run("medium list with paiganation", func(t *testing.T) {

		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(mediumlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(2, time.Now(), time.Now(), nil, mediumlist[1]["name"], mediumlist[1]["slug"], mediumlist[1]["type"], mediumlist[1]["title"], mediumlist[1]["description"], mediumlist[1]["caption"], mediumlist[1]["alt_text"], mediumlist[1]["file_size"], mediumlist[1]["url"], mediumlist[1]["dimensions"]))

		e.GET("/media").
			WithQueryObject(map[string]interface{}{
				"limit": 1,
				"page":  2,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(mediumlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(mediumlist[1])

		mock.ExpectationsWereMet()

	})

}
