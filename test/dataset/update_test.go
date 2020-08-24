package dataset

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
)

func TestUpdateDataset(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("update dataset", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(DatasetCols).
				AddRow(1, time.Now(), time.Now(), nil, "title", "description", "source", "frequency", "temporal_coverage", "granularity", "contact_name", "contact_email", "license", "data_standard", nilJsonb(), 10, 200, 2, 2))

		tagAssociationSelectMock(mock)

		mock.ExpectBegin()

		tag.TagSelectMock(mock)

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_dataset_tag"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(`UPDATE \"dp_dataset\" SET (.+)  WHERE (.+) \"dp_dataset\".\"id\" = `).
			WithArgs(Dataset["contact_email"], Dataset["contact_name"], Dataset["currency_id"], Dataset["data_standard"], Dataset["description"], Dataset["featured_medium_id"], Dataset["frequency"], Dataset["granularity"], Dataset["license"], Dataset["price"], Dataset["related_articles"], Dataset["source"], Dataset["temporal_coverage"], Dataset["time_saved"], Dataset["title"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`INSERT INTO "dp_dataset_tag"`).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		DatasetSelectMock(mock)

		medium.MediumSelectMock(mock)

		currency.CurrencySelectMock(mock)

		tag.TagSelectMock(mock)

		datasetFormatSelectMock(mock, 1)

		mock.ExpectCommit()

		result := e.PUT(path).
			WithPath("dataset_id", "1").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(DatasetReceive)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("dataset record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(DatasetCols))

		e.PUT(path).
			WithPath("dataset_id", "1").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable dataset body", func(t *testing.T) {
		e.PUT(path).
			WithPath("dataset_id", "1").
			WithJSON(invalidDataset).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid dataset id", func(t *testing.T) {
		e.PUT(path).
			WithPath("dataset_id", "abc").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("new featured medium does not exist", func(t *testing.T) {
		updateWithErrorMock(mock, errDatasetMediumFK)

		e.PUT(path).
			WithPath("dataset_id", "1").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("new currency does not exist", func(t *testing.T) {
		updateWithErrorMock(mock, errDatasetCurrencyFK)

		e.PUT(path).
			WithPath("dataset_id", "1").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
