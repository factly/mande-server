package datasetformat

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

func TestDeleteDatasetFormat(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("delete dataset format", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows(DatasetFormatCols).
				AddRow(1, time.Now(), time.Now(), nil, DatasetFormat["format_id"], 1, DatasetFormat["url"]))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_dataset_format" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(path).
			WithPathObject(map[string]interface{}{
				"dataset_id": "1",
				"format_id":  "1",
			}).
			Expect().
			Status(http.StatusOK)

		test.ExpectationsMet(t, mock)
	})

	t.Run("dataset format record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows(DatasetFormatCols))

		e.DELETE(path).
			WithPathObject(map[string]interface{}{
				"dataset_id": "1",
				"format_id":  "1",
			}).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})
}
