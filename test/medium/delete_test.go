package medium

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestDeleteMedium(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	test.KavachGock()
	test.KetoGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("delete medium", func(t *testing.T) {
		MediumSelectMock(mock, 1)

		mediumCatalogExpect(mock, 0)

		mediumDatasetExpect(mock, 0)

		mediumProductExpect(mock, 0)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_medium" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("media_id", "1").
			Expect().
			Status(http.StatusOK)

		test.ExpectationsMet(t, mock)
	})

	t.Run("medium record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(MediumCols))

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("media_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid medium id", func(t *testing.T) {
		e.DELETE(path).
			WithHeaders(headers).
			WithPath("media_id", "abc").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("medium is associated with catalog", func(t *testing.T) {
		MediumSelectMock(mock, 1)

		mediumCatalogExpect(mock, 1)

		e.DELETE(path).
			WithPath("media_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("medium is associated with dataset", func(t *testing.T) {
		MediumSelectMock(mock, 1)

		mediumCatalogExpect(mock, 0)

		mediumDatasetExpect(mock, 1)

		e.DELETE(path).
			WithPath("media_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("medium is associated with product", func(t *testing.T) {
		MediumSelectMock(mock, 1)

		mediumCatalogExpect(mock, 0)

		mediumDatasetExpect(mock, 0)

		mediumProductExpect(mock, 1)

		e.DELETE(path).
			WithPath("media_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("delete medium when meili is down", func(t *testing.T) {
		gock.Off()
		test.KavachGock()
		test.KetoGock()
		gock.New(server.URL).EnableNetworking().Persist()

		MediumSelectMock(mock, 1)

		mediumCatalogExpect(mock, 0)

		mediumDatasetExpect(mock, 0)

		mediumProductExpect(mock, 0)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_medium" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectRollback()

		e.DELETE(path).
			WithPath("media_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
