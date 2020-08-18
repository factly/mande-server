package tag

import (
	"database/sql/driver"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/model"
	"github.com/gavv/httpexpect"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func SetupTestDB() sqlmock.Sqlmock {
	DSN := "sqlmock_db_0"
	_, mock, err := sqlmock.NewWithDSN(DSN)
	if err != nil {
		fmt.Println(err)
	}

	os.Setenv("DSN", DSN)

	model.SetupDB()

	// model.DB, err = gorm.Open("postgres", db)

	// model.DB.LogMode(true)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	return mock
}

func TestGetTagList(t *testing.T) {

	// Setup DB
	// model.SetupDB()
	mock := SetupTestDB()

	// gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	// 	return "dp_" + defaultTableName
	// }

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	// DB
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_tags"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tags"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}))

	// Request
	e.GET("/tags").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsMap(map[string]interface{}{"total": 0})

	//CREATE TAG
	createdTag := map[string]interface{}{
		"title": "Test Tag",
		"slug":  "test-tag",
	}

	// DB
	const sqlInsert = `
	INSERT INTO "dp_tags"
	`
	mock.ExpectBegin()
	mock.ExpectQuery(sqlInsert).
		WithArgs(AnyTime{}, AnyTime{}, nil, createdTag["title"], createdTag["slug"]).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tags"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}).
			AddRow(1, time.Now(), time.Now(), nil, createdTag["title"], createdTag["slug"]))

	e.POST("/tags").
		WithJSON(createdTag).
		Expect().
		Status(http.StatusCreated).
		JSON().
		Object().
		ContainsMap(createdTag)

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
