package cart

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/factly/data-portal-server/test/membership"
	"github.com/factly/data-portal-server/test/plan"
	"github.com/factly/data-portal-server/test/product"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestDetailCart(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	// ADMIN tests
	CommonDetailTests(t, mock, adminExpect)

	t.Run("get cart item by id", func(t *testing.T) {
		CartItemSelectMock(mock, 1)

		membership.MembershipSelectMock(mock)
		plan.PlanSelectMock(mock)
		product.ProductSelectMock(mock)
		currency.CurrencySelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product_dataset"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"product_id", "dataset_id"}).
				AddRow(1, 1))
		dataset.DatasetSelectMock(mock)

		medium.MediumSelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product_tag"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"product_id", "tag_id"}).
				AddRow(1, 1))

		tag.TagSelectMock(mock)

		adminExpect.GET(path).
			WithHeaders(headers).
			WithPath("cartitem_id", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(CartItem)

		test.ExpectationsMet(t, mock)
	})

	t.Run("cart item record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CartItemCols))

		adminExpect.GET(path).
			WithHeaders(headers).
			WithPath("cartitem_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	// USER test
	CommonDetailTests(t, mock, userExpect)

	t.Run("get cart item by id", func(t *testing.T) {
		CartItemSelectMock(mock, 1, 1)

		membership.MembershipSelectMock(mock)
		plan.PlanSelectMock(mock)
		product.ProductSelectMock(mock)
		currency.CurrencySelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product_dataset"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"product_id", "dataset_id"}).
				AddRow(1, 1))
		dataset.DatasetSelectMock(mock)

		medium.MediumSelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product_tag"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"product_id", "tag_id"}).
				AddRow(1, 1))
		tag.TagSelectMock(mock)

		userExpect.GET(path).
			WithHeaders(headers).
			WithPath("cartitem_id", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(CartItem)

		test.ExpectationsMet(t, mock)
	})

	t.Run("cart item record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows(CartItemCols))

		userExpect.GET(path).
			WithHeaders(headers).
			WithPath("cartitem_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid user header", func(t *testing.T) {
		userExpect.GET(path).
			WithHeaders(headers).
			WithPath("cartitem_id", "1").
			Expect().
			Status(http.StatusNotFound)
	})

	server.Close()
}

func CommonDetailTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {

	t.Run("invalid cart item id", func(t *testing.T) {
		e.GET(path).
			WithHeaders(headers).
			WithPath("cartitem_id", "abc").
			Expect().
			Status(http.StatusBadRequest)
	})
}
