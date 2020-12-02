package dataset

import (
	"errors"
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
	test.KetoGock()
	test.KavachGock()
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
			WithHeaders(headers).
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
			WithHeaders(headers).
			WithPath("dataset_id", "1").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable dataset body", func(t *testing.T) {
		e.PUT(path).
			WithHeaders(headers).
			WithPath("dataset_id", "1").
			WithJSON(invalidDataset).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("undecodable dataset body", func(t *testing.T) {
		e.PUT(path).
			WithHeaders(headers).
			WithPath("dataset_id", "1").
			WithJSON(undecodableDataset).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid dataset id", func(t *testing.T) {
		e.PUT(path).
			WithHeaders(headers).
			WithPath("dataset_id", "abc").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusBadRequest)
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
			WithHeaders(headers).
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
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(DatasetCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, "title", "description", "source", "frequency", "temporal_coverage", "granularity", "contact_name", "contact_email", "license", "data_standard", "sample_url", nilJsonb(), 10, 200, 2, 2))

		mock.ExpectBegin()

		tag.TagSelectMock(mock)

		mock.ExpectQuery(`INSERT INTO "dp_tag"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, tag.Tag["title"], tag.Tag["slug"], 1).
			WillReturnError(errors.New(`cannot replace tags`))

		mock.ExpectExec(`INSERT INTO "dp_dataset_tag"`).
			WithArgs(1, 1).
			WillReturnError(errors.New(`cannot replace tags`))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_dataset_tag"`)).
			WillReturnError(errors.New(`cannot replace tags`))

		mock.ExpectRollback()

		e.PUT(path).
			WithPath("dataset_id", "1").
			WithHeaders(headers).
			WithJSON(Dataset).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

	t.Run("new featured medium does not exist", func(t *testing.T) {
		updateMock(mock, errDatasetMediumFK)

		e.PUT(path).
			WithHeaders(headers).
			WithPath("dataset_id", "1").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("new currency does not exist", func(t *testing.T) {
		updateMock(mock, errDatasetCurrencyFK)

		e.PUT(path).
			WithHeaders(headers).
			WithPath("dataset_id", "1").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("update dataset when meili is down", func(t *testing.T) {
		gock.Off()
		test.KetoGock()
		test.KavachGock()
		gock.New(server.URL).EnableNetworking().Persist()

		updateMock(mock, nil)

		DatasetSelectMock(mock)

		currency.CurrencySelectMock(mock)
		medium.MediumSelectMock(mock)

		tagAssociationSelectMock(mock)
		datasetFormatSelectMock(mock, 1)

		mock.ExpectRollback()

		e.PUT(path).
			WithHeaders(headers).
			WithPath("dataset_id", "1").
			WithJSON(Dataset).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
