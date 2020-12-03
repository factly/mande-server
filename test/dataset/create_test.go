package dataset

import (
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

func TestCreateDataset(t *testing.T) {
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

	t.Run("create a dataset", func(t *testing.T) {
		tag.TagSelectMock(mock)
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_dataset"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Dataset["title"], Dataset["description"], Dataset["source"], Dataset["frequency"], Dataset["temporal_coverage"], Dataset["granularity"], Dataset["contact_name"], Dataset["contact_email"], Dataset["license"], Dataset["data_standard"], Dataset["sample_url"], Dataset["related_articles"], Dataset["time_saved"], Dataset["price"], Dataset["currency_id"], Dataset["featured_medium_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"featured_medium_id", "id"}).AddRow(1, 1))

		mock.ExpectQuery(`INSERT INTO "dp_tag"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, tag.Tag["title"], tag.Tag["slug"], 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(`INSERT INTO "dp_dataset_tag"`).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		DatasetSelectMock(mock)

		currency.CurrencySelectMock(mock)
		medium.MediumSelectMock(mock)

		tagAssociationSelectMock(mock)

		mock.ExpectCommit()

		result := e.POST(basePath).
			WithHeaders(headers).
			WithJSON(Dataset).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(DatasetReceive)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable dataset body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidDataset).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty dataset body", func(t *testing.T) {
		e.POST(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("featured medium does not exist", func(t *testing.T) {
		insertWithErrorMock(mock, errDatasetMediumFK)

		e.POST(basePath).
			WithJSON(Dataset).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("currency does not exist", func(t *testing.T) {
		insertWithErrorMock(mock, errDatasetCurrencyFK)

		e.POST(basePath).
			WithJSON(Dataset).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a dataset when meili is down", func(t *testing.T) {
		gock.Off()
		test.KetoGock()
		test.KavachGock()
		gock.New(server.URL).EnableNetworking().Persist()

		tag.TagSelectMock(mock)
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_dataset"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Dataset["title"], Dataset["description"], Dataset["source"], Dataset["frequency"], Dataset["temporal_coverage"], Dataset["granularity"], Dataset["contact_name"], Dataset["contact_email"], Dataset["license"], Dataset["data_standard"], Dataset["sample_url"], Dataset["related_articles"], Dataset["time_saved"], Dataset["price"], Dataset["currency_id"], Dataset["featured_medium_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"featured_medium_id", "id"}).AddRow(1, 1))

		mock.ExpectQuery(`INSERT INTO "dp_tag"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, tag.Tag["title"], tag.Tag["slug"], 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(`INSERT INTO "dp_dataset_tag"`).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		DatasetSelectMock(mock)

		currency.CurrencySelectMock(mock)
		medium.MediumSelectMock(mock)

		tagAssociationSelectMock(mock)

		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Dataset).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
