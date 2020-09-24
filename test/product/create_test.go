package product

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestCreateProduct(t *testing.T) {
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

	t.Run("create a product", func(t *testing.T) {
		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_product"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Product["title"], Product["slug"], Product["price"], Product["status"], Product["currency_id"], Product["featured_medium_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(`INSERT INTO "dp_product_tag"`).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(`INSERT INTO "dp_product_dataset"`).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tagsAssociationSelectMock(mock, 1)

		datasetsAssociationSelectMock(mock, 1)

		mock.ExpectCommit()

		result := e.POST(basePath).
			WithJSON(Product).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(ProductReceive)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable product body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidProduct).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty product body", func(t *testing.T) {
		e.POST(basePath).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("featured medium does not exist", func(t *testing.T) {
		insertWithErrorMock(mock, errProductMediumFK)

		e.POST(basePath).
			WithJSON(Product).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("currency does not exist", func(t *testing.T) {
		insertWithErrorMock(mock, errProductCurrencyFK)

		e.POST(basePath).
			WithJSON(Product).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a product when meili is down", func(t *testing.T) {
		gock.Off()
		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_product"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Product["title"], Product["slug"], Product["price"], Product["status"], Product["currency_id"], Product["featured_medium_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(`INSERT INTO "dp_product_tag"`).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(`INSERT INTO "dp_product_dataset"`).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tagsAssociationSelectMock(mock, 1)

		datasetsAssociationSelectMock(mock, 1)

		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Product).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
