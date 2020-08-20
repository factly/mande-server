package tag

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
)

func TestDeleteTag(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("delete tag", func(t *testing.T) {
		tagSelectMock(mock)

		tagProductExpect(mock, 0)

		tagDatasetExpect(mock, 0)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_tag" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(path).
			WithPath("tag_id", "1").
			Expect().
			Status(http.StatusOK)

		mock.ExpectationsWereMet()
	})

	t.Run("tag not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(tagCols))

		e.DELETE(path).
			WithPath("tag_id", "1").
			Expect().
			Status(http.StatusNotFound)

		mock.ExpectationsWereMet()
	})

	t.Run("invalid tag id", func(t *testing.T) {
		e.DELETE(path).
			WithPath("tag_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("tag associated with product", func(t *testing.T) {
		tagSelectMock(mock)

		tagProductExpect(mock, 1)

		e.DELETE(path).
			WithPath("tag_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		mock.ExpectationsWereMet()
	})

	t.Run("tag associated with dataset", func(t *testing.T) {
		tagSelectMock(mock)

		tagProductExpect(mock, 0)

		tagDatasetExpect(mock, 1)

		e.DELETE(path).
			WithPath("tag_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		mock.ExpectationsWereMet()
	})
}
