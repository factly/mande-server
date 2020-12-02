package cart

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/factly/data-portal-server/test/plan"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/factly/data-portal-server/test/membership"
	"github.com/factly/data-portal-server/test/product"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateCart(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterUserRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("create a cart item", func(t *testing.T) {
		mock.ExpectBegin()

		membership.MembershipSelectMock(mock)

		mock.ExpectQuery(`INSERT INTO "dp_cart_item"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, CartItem["status"], 1, CartItem["product_id"], CartItem["membership_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id", "membership_id"}).AddRow(1, 1))

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

		mock.ExpectCommit()

		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(CartItem).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(CartItem)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable cart item body", func(t *testing.T) {
		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(invalidCartItem).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("undecodable cart item body", func(t *testing.T) {
		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(undecodableCartItem).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty cart item body", func(t *testing.T) {
		e.POST(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("product does not exist", func(t *testing.T) {

		mock.ExpectBegin()
		membership.MembershipSelectMock(mock)
		mock.ExpectQuery(`INSERT INTO "dp_cart_item"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, CartItem["status"], 1, CartItem["product_id"], CartItem["membership_id"]).
			WillReturnError(errCartItemProductFK)
		mock.ExpectRollback()

		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(CartItem).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("membership does not exist", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_membership"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(membership.MembershipCols))

		mock.ExpectRollback()

		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(CartItem).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a cart item when meili is down", func(t *testing.T) {
		gock.Off()

		mock.ExpectBegin()
		membership.MembershipSelectMock(mock)

		mock.ExpectQuery(`INSERT INTO "dp_cart_item"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, CartItem["status"], 1, CartItem["product_id"], CartItem["membership_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id", "membership_id"}).AddRow(1, 1))

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

		mock.ExpectRollback()

		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(CartItem).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})
}
