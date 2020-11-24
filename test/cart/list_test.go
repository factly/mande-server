package cart

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
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/factly/data-portal-server/test/membership"
	"github.com/factly/data-portal-server/test/plan"
	"github.com/factly/data-portal-server/test/product"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestListCartItems(t *testing.T) {

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

	// ADMIN specific tests
	CommonListTests(t, mock, adminExpect)

	t.Run("get cart item list with user query", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(cartitemslist)))

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CartItemCols).
				AddRow(1, time.Now(), time.Now(), nil, cartitemslist[0]["status"], cartitemslist[0]["user_id"], cartitemslist[0]["product_id"], cartitemslist[0]["membership_id"]).
				AddRow(2, time.Now(), time.Now(), nil, cartitemslist[1]["status"], cartitemslist[1]["user_id"], cartitemslist[1]["product_id"], cartitemslist[1]["membership_id"]))

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

		delete(cartitemslist[0], "product_id")

		adminExpect.GET(basePath).
			WithHeaders(headers).
			WithQuery("user", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(cartitemslist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(cartitemslist[0])

		cartitemslist[0]["product_id"] = 1

		test.ExpectationsMet(t, mock)
	})

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	// USER specific tests
	CommonListTests(t, mock, userExpect)

	t.Run("invalid user header", func(t *testing.T) {
		userExpect.GET(basePath).
			WithHeaders(map[string]string{
				"X-Organisation": "1",
				"X-User":         "invalid",
			}).
			Expect().
			Status(http.StatusUnauthorized)
	})

	server.Close()
}

func CommonListTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get empty cart list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(CartItemCols))

		e.GET(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("get cart item list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(cartitemslist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(CartItemCols).
				AddRow(1, time.Now(), time.Now(), nil, cartitemslist[0]["status"], cartitemslist[0]["user_id"], cartitemslist[0]["product_id"], cartitemslist[0]["membership_id"]).
				AddRow(2, time.Now(), time.Now(), nil, cartitemslist[1]["status"], cartitemslist[1]["user_id"], cartitemslist[1]["product_id"], cartitemslist[1]["membership_id"]))

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

		delete(cartitemslist[0], "product_id")

		e.GET(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(cartitemslist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(cartitemslist[0])

		cartitemslist[0]["product_id"] = 1

		test.ExpectationsMet(t, mock)
	})

	t.Run("get cart item list with paiganation", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(cartitemslist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(CartItemCols).
				AddRow(2, time.Now(), time.Now(), nil, cartitemslist[1]["status"], cartitemslist[1]["user_id"], cartitemslist[1]["product_id"], cartitemslist[1]["membership_id"]))

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

		delete(cartitemslist[1], "product_id")

		e.GET(basePath).
			WithHeaders(headers).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(cartitemslist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(cartitemslist[1])

		cartitemslist[1]["product_id"] = 1

		test.ExpectationsMet(t, mock)

	})
}
