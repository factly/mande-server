package dataset

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestUpdateDataset(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("update dataset", func(t *testing.T) {
		updateMock(mock, nil)

		DatasetSelectMock(mock)

		currency.CurrencySelectMock(mock)
		medium.MediumSelectMock(mock)

		tagAssociationSelectMock(mock)
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

	t.Run("undecodable dataset body", func(t *testing.T) {
		e.PUT(path).
			WithPath("dataset_id", "1").
			WithJSON(undecodableDataset).
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

	t.Run("update dataset with featured_medium_id = 0", func(t *testing.T) {
		updateWithoutFeaturedMedium(mock)

		DatasetSelectMock(mock)

		currency.CurrencySelectMock(mock)
		medium.MediumSelectMock(mock)

		tagAssociationSelectMock(mock)
		datasetFormatSelectMock(mock, 1)

		mock.ExpectCommit()

		Dataset["featured_medium_id"] = 0
		result := e.PUT(path).
			WithPath("dataset_id", "1").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(DatasetReceive)

		validateAssociations(result)
		Dataset["featured_medium_id"] = 1
		test.ExpectationsMet(t, mock)
	})

	t.Run("replacing old tags fail", func(t *testing.T) {
		preUpdateMock(mock)

		mock.ExpectExec(`UPDATE \"dp_dataset\"`).
			WithArgs(test.AnyTime{}, Dataset["title"], Dataset["description"], Dataset["source"], Dataset["frequency"], Dataset["temporal_coverage"], Dataset["granularity"], Dataset["contact_name"], Dataset["contact_email"], Dataset["license"], Dataset["data_standard"], Dataset["sample_url"], Dataset["related_articles"], Dataset["time_saved"], Dataset["price"], Dataset["currency_id"], Dataset["featured_medium_id"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`INSERT INTO "dp_tag"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, tag.Tag["title"], tag.Tag["slug"], 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(`INSERT INTO "dp_dataset_tag"`).
			WithArgs(1, 1).
			WillReturnError(errors.New(`cannot replace tags`))
		mock.ExpectRollback()

		e.PUT(path).
			WithPath("dataset_id", "1").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

	t.Run("new featured medium does not exist", func(t *testing.T) {
		updateMock(mock, errDatasetMediumFK)

		e.PUT(path).
			WithPath("dataset_id", "1").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("new currency does not exist", func(t *testing.T) {
		updateMock(mock, errDatasetCurrencyFK)

		e.PUT(path).
			WithPath("dataset_id", "1").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("update dataset when meili is down", func(t *testing.T) {
		gock.Off()
		updateMock(mock, nil)

		DatasetSelectMock(mock)

		currency.CurrencySelectMock(mock)
		medium.MediumSelectMock(mock)

		tagAssociationSelectMock(mock)
		datasetFormatSelectMock(mock, 1)

		mock.ExpectRollback()

		e.PUT(path).
			WithPath("dataset_id", "1").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
