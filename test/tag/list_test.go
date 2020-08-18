package tag

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/test"

	"github.com/factly/data-portal-server/action"
	"github.com/gavv/httpexpect"
)

func TestGetTagsList(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	// DB
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_tag"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}))

	// Request
	e.GET("/tags").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsMap(map[string]interface{}{"total": 0})

	mock.ExpectationsWereMet()

	// id := strconv.Itoa(int(res.Value("id").Number().Raw()))

	// e.GET("/tags/" + id).
	// 	Expect().
	// 	Status(http.StatusOK).
	// 	JSON().
	// 	Object().
	// 	ContainsMap(createdTag)

	// updatedTag := map[string]interface{}{
	// 	"title": "Test Tag Updated",
	// 	"slug":  "test-tag-updated",
	// }

	// e.PUT("/tags/" + id).
	// 	WithJSON(updatedTag).
	// 	Expect().
	// 	Status(http.StatusOK).
	// 	JSON().
	// 	Object().
	// 	ContainsMap(updatedTag)

	// model.DB.Delete(&model.Tag{})
}
