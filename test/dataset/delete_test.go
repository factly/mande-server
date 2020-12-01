package dataset

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestDeleteDataset(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("delete dataset", func(t *testing.T) {
		DatasetSelectMock(mock)

		tagAssociationSelectMock(mock)

		datasetProductExpect(mock, 0)

		deleteMock(mock, nil)
		mock.ExpectCommit()

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("dataset_id", "1").
			Expect().
			Status(http.StatusOK)

		test.ExpectationsMet(t, mock)
	})

	t.Run("dataset record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(DatasetCols))

		e.DELETE(path).
			WithPath("dataset_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid dataset id", func(t *testing.T) {
		e.DELETE(path).
			WithHeaders(headers).
			WithPath("dataset_id", "abc").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("deleting dataset fails", func(t *testing.T) {
		DatasetSelectMock(mock)

		tagAssociationSelectMock(mock)

		datasetProductExpect(mock, 0)

		deleteMock(mock, errors.New("cannot delete dataset"))
		mock.ExpectRollback()

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("dataset_id", "1").
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("dataset associated with product", func(t *testing.T) {
		DatasetSelectMock(mock)

		tagAssociationSelectMock(mock)

		datasetProductExpect(mock, 1)

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("dataset_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("delete dataset when meili is down", func(t *testing.T) {
		gock.Off()
		test.KetoGock()
		test.KavachGock()
		gock.New(server.URL).EnableNetworking().Persist()

		DatasetSelectMock(mock)

		tagAssociationSelectMock(mock)

		datasetProductExpect(mock, 0)

		deleteMock(mock, nil)
		mock.ExpectRollback()

		e.DELETE(path).
			WithPath("dataset_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
