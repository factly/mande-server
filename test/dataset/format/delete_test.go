package datasetformat

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestDeleteDatasetFormat(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()

	e := httpexpect.New(t, server.URL)

	t.Run("delete dataset format", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows(DatasetFormatCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, DatasetFormat["format_id"], 1, DatasetFormat["url"]))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_dataset_format" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(path).
			WithHeaders(headers).
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
			WithHeaders(headers).
			WithPathObject(map[string]interface{}{
				"dataset_id": "1",
				"format_id":  "1",
			}).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid dataset id", func(t *testing.T) {
		e.DELETE(path).
			WithHeaders(headers).
			WithPathObject(map[string]interface{}{
				"dataset_id": "abc",
				"format_id":  "1",
			}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("invalid format id", func(t *testing.T) {
		e.DELETE(path).
			WithHeaders(headers).
			WithPathObject(map[string]interface{}{
				"dataset_id": "1",
				"format_id":  "abc",
			}).
			Expect().
			Status(http.StatusBadRequest)
	})
}
