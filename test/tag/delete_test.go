package tag

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestDeleteTag(t *testing.T) {

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

	t.Run("delete tag", func(t *testing.T) {
		TagSelectMock(mock)

		tagProductExpect(mock, 0)

		tagDatasetExpect(mock, 0)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_tag" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(path).
			WithHeaders(headers).
			WithPath("tag_id", "1").
			Expect().
			Status(http.StatusOK)

		test.ExpectationsMet(t, mock)
	})

	t.Run("tag not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(TagCols))

		e.DELETE(path).
			WithPath("tag_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid tag id", func(t *testing.T) {
		e.DELETE(path).
			WithPath("tag_id", "abc").
			WithHeaders(headers).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("tag associated with product", func(t *testing.T) {
		TagSelectMock(mock)

		tagProductExpect(mock, 1)

		e.DELETE(path).
			WithPath("tag_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("tag associated with dataset", func(t *testing.T) {
		TagSelectMock(mock)

		tagProductExpect(mock, 0)

		tagDatasetExpect(mock, 1)

		e.DELETE(path).
			WithPath("tag_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("delete tag when meili is down", func(t *testing.T) {
		gock.Off()
		test.KavachGock()
		test.KetoGock()
		gock.New(server.URL).EnableNetworking().Persist()

		TagSelectMock(mock)

		tagProductExpect(mock, 0)

		tagDatasetExpect(mock, 0)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_tag" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectRollback()

		e.DELETE(path).
			WithPath("tag_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
