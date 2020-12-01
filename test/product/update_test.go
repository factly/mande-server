package product

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
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestUpdateProduct(t *testing.T) {
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

	t.Run("update product", func(t *testing.T) {
		updateMock(mock, nil)

		ProductSelectMock(mock, 1, 1)

		currency.CurrencySelectMock(mock)
		datasetsAssociationSelectMock(mock, 1)
		medium.MediumSelectMock(mock)
		tagsAssociationSelectMock(mock, 1)

		mock.ExpectCommit()

		result := e.PUT(path).
			WithPath("product_id", "1").
			WithHeaders(headers).
			WithJSON(Product).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(ProductReceive)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("product record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(ProductCols))

		e.PUT(path).
			WithPath("product_id", "1").
			WithHeaders(headers).
			WithJSON(Product).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable product body", func(t *testing.T) {
		e.PUT(path).
			WithPath("product_id", "1").
			WithHeaders(headers).
			WithJSON(invalidProduct).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("undecodable product body", func(t *testing.T) {
		e.PUT(path).
			WithPath("product_id", "1").
			WithHeaders(headers).
			WithJSON(undecodableProduct).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid product id", func(t *testing.T) {
		e.PUT(path).
			WithPath("product_id", "abc").
			WithHeaders(headers).
			WithJSON(Product).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("replacing old tags fails", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(ProductCols).
				AddRow(1, time.Now(), time.Now(), nil, "title", "slug", 200, "status", 2, 2))

		mock.ExpectBegin()

		tag.TagSelectMock(mock)

		mock.ExpectQuery(`INSERT INTO "dp_tag"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, tag.Tag["title"], tag.Tag["slug"], 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(`INSERT INTO "dp_product_tag"`).
			WithArgs(1, 1).
			WillReturnError(errors.New(`cannot replace tags`))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_tag"`)).
			WillReturnError(errors.New(`cannot replace tags`))

		mock.ExpectRollback()

		e.PUT(path).
			WithPath("product_id", "1").
			WithHeaders(headers).
			WithJSON(Product).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

	t.Run("replacing old datasets fails", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(ProductCols).
				AddRow(1, time.Now(), time.Now(), nil, "title", "slug", 200, "status", 2, 2))

		mock.ExpectBegin()

		tag.TagSelectMock(mock)

		mock.ExpectQuery(`INSERT INTO "dp_tag"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, tag.Tag["title"], tag.Tag["slug"], 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(`INSERT INTO "dp_product_tag"`).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_tag"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		dataset.DatasetSelectMock(mock)

		mock.ExpectQuery(`INSERT INTO "dp_dataset"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, dataset.Dataset["title"], dataset.Dataset["description"], dataset.Dataset["source"], dataset.Dataset["frequency"], dataset.Dataset["temporal_coverage"], dataset.Dataset["granularity"], dataset.Dataset["contact_name"], dataset.Dataset["contact_email"], dataset.Dataset["license"], dataset.Dataset["data_standard"], dataset.Dataset["sample_url"], dataset.Dataset["related_articles"], dataset.Dataset["time_saved"], dataset.Dataset["price"], dataset.Dataset["currency_id"], dataset.Dataset["featured_medium_id"], 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(`INSERT INTO "dp_product_dataset"`).
			WithArgs(1, 1).
			WillReturnError(errors.New(`cannot replace datasets`))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_dataset"`)).
			WillReturnError(errors.New(`cannot replace datasets`))

		mock.ExpectRollback()

		e.PUT(path).
			WithHeaders(headers).
			WithPath("product_id", "1").
			WithJSON(Product).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

	t.Run("update product with new featured_medium_id = 0", func(t *testing.T) {
		updateMockWithoutMedium(mock)

		ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)
		datasetsAssociationSelectMock(mock, 1)
		medium.MediumSelectMock(mock)
		tagsAssociationSelectMock(mock, 1)
		mock.ExpectCommit()

		Product["featured_medium_id"] = 0
		result := e.PUT(path).
			WithPath("product_id", "1").
			WithJSON(Product).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(ProductReceive)

		validateAssociations(result)

		Product["featured_medium_id"] = 1

		test.ExpectationsMet(t, mock)
	})

	t.Run("new featured medium does not exist", func(t *testing.T) {
		updateMock(mock, errProductMediumFK)

		e.PUT(path).
			WithPath("product_id", "1").
			WithJSON(Product).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("new currency does not exist", func(t *testing.T) {
		updateMock(mock, errProductCurrencyFK)

		e.PUT(path).
			WithPath("product_id", "1").
			WithJSON(Product).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("update product when meili is down", func(t *testing.T) {
		gock.Off()
		test.KavachGock()
		test.KetoGock()
		gock.New(server.URL).EnableNetworking().Persist()

		updateMock(mock, nil)

		ProductSelectMock(mock, 1, 1)

		currency.CurrencySelectMock(mock)
		datasetsAssociationSelectMock(mock, 1)
		medium.MediumSelectMock(mock)
		tagsAssociationSelectMock(mock, 1)

		mock.ExpectRollback()

		e.PUT(path).
			WithPath("product_id", "1").
			WithHeaders(headers).
			WithJSON(Product).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

}
