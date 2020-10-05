package cart

import (
	"net/http"
	"net/http/httptest"
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
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("create a cart item", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_cart_item"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, CartItem["status"], 1, CartItem["product_id"], CartItem["membership_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		CartItemSelectMock(mock)

		product.ProductSelectMock(mock)

		membership.MembershipSelectMock(mock)

		plan.PlanSelectMock(mock)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		mock.ExpectCommit()

		e.POST(basePath).
			WithHeader("X-User", "1").
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
			WithHeader("X-User", "1").
			WithJSON(invalidCartItem).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("undecodable cart item body", func(t *testing.T) {
		e.POST(basePath).
			WithHeader("X-User", "1").
			WithJSON(undecodableCartItem).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty cart item body", func(t *testing.T) {
		e.POST(basePath).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid user header", func(t *testing.T) {
		e.POST(basePath).
			WithHeader("X-User", "abc").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("user does not exist", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_cart_item"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, CartItem["status"], 1, CartItem["product_id"], CartItem["membership_id"]).
			WillReturnError(errCartItemProductFK)
		mock.ExpectRollback()

		e.POST(basePath).
			WithHeader("X-User", "1").
			WithJSON(CartItem).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a cart item when meili is down", func(t *testing.T) {
		gock.Off()

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_cart_item"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, CartItem["status"], 1, CartItem["product_id"], CartItem["membership_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		CartItemSelectMock(mock)

		product.ProductSelectMock(mock)

		membership.MembershipSelectMock(mock)

		plan.PlanSelectMock(mock)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		mock.ExpectRollback()

		e.POST(basePath).
			WithHeader("X-User", "1").
			WithJSON(CartItem).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})
}
