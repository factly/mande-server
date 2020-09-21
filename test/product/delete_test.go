package product

import (
	"errors"
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

func TestDeleteProduct(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("delete product", func(t *testing.T) {
		ProductSelectMock(mock)

		tagsAssociationSelectMock(mock, 1)

		datasetsAssociationSelectMock(mock, 1)

		productCatalogExpect(mock, 0)

		productOrderExpect(mock, 0)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_tag"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_dataset"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_product" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(path).
			WithPath("product_id", "1").
			Expect().
			Status(http.StatusOK)

		test.ExpectationsMet(t, mock)
	})

	t.Run("product record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(ProductCols))

		e.DELETE(path).
			WithPath("product_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid product id", func(t *testing.T) {
		e.DELETE(path).
			WithPath("product_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("deleting product fails", func(t *testing.T) {
		ProductSelectMock(mock)

		tagsAssociationSelectMock(mock, 1)

		datasetsAssociationSelectMock(mock, 1)

		productCatalogExpect(mock, 0)

		productOrderExpect(mock, 0)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_tag"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_dataset"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_product" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnError(errors.New("cannot delete"))
		mock.ExpectRollback()

		e.DELETE(path).
			WithPath("product_id", "1").
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("product is associated with catalog", func(t *testing.T) {
		ProductSelectMock(mock)

		tagsAssociationSelectMock(mock, 1)

		datasetsAssociationSelectMock(mock, 1)

		productCatalogExpect(mock, 1)

		e.DELETE(path).
			WithPath("product_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("product is associated with order", func(t *testing.T) {
		ProductSelectMock(mock)

		tagsAssociationSelectMock(mock, 1)

		datasetsAssociationSelectMock(mock, 1)

		productCatalogExpect(mock, 0)

		productOrderExpect(mock, 1)

		e.DELETE(path).
			WithPath("product_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("delete product when meili is down", func(t *testing.T) {
		gock.Off()
		ProductSelectMock(mock)

		tagsAssociationSelectMock(mock, 1)

		datasetsAssociationSelectMock(mock, 1)

		productCatalogExpect(mock, 0)

		productOrderExpect(mock, 0)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_tag"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_dataset"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_product" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectRollback()

		e.DELETE(path).
			WithPath("product_id", "1").
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
