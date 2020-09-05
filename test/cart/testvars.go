package cart

import (
	"errors"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/product"
)

var Cart map[string]interface{} = map[string]interface{}{
	"status":      "teststatus",
	"user_id":     1,
	"product_ids": []uint{1},
}

var CartReceive map[string]interface{} = map[string]interface{}{
	"status":  "teststatus",
	"user_id": 1,
}

var invalidCart map[string]interface{} = map[string]interface{}{
	"stat":      "testtus",
	"userid":    1,
	"productds": []uint{1},
}

var undecodableCart map[string]interface{} = map[string]interface{}{
	"status":  45,
	"user_id": "1",
}

var cartlist []map[string]interface{} = []map[string]interface{}{
	{
		"status":      "teststatus1",
		"user_id":     1,
		"product_ids": []uint{1},
	},
	{
		"status":      "teststatus2",
		"user_id":     1,
		"product_ids": []uint{1},
	},
}

var CartCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "status", "user_id"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_cart"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_cart"`)
var errCartProductFK error = errors.New(`pq: insert or update on table "dp_cart" violates foreign key constraint "dp_cart_user_id_dp_user_id_foreign"`)

const basePath string = "/carts"
const path string = "/carts/{cart_id}"

func CartSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(CartCols).
			AddRow(1, time.Now(), time.Now(), nil, Cart["status"], Cart["user_id"]))
}

func productsAssociationSelectMock(mock sqlmock.Sqlmock, cartId int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product" INNER JOIN "dp_cart_item"`)).
		WithArgs(cartId).
		WillReturnRows(sqlmock.NewRows(append(product.ProductCols, []string{"product_id", "cart_id"}...)).
			AddRow(1, time.Now(), time.Now(), nil, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1, cartId))
}

func cartOrderExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_order`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func preUpdateMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(CartCols).
			AddRow(1, time.Now(), time.Now(), nil, "status", 2))

	productsAssociationSelectMock(mock, 1)

	mock.ExpectBegin()

	product.ProductSelectMock(mock)

}

func updateMock(mock sqlmock.Sqlmock, err error) {
	preUpdateMock(mock)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_cart_item"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err != nil {
		mock.ExpectExec(`UPDATE \"dp_cart\" SET (.+)  WHERE (.+) \"dp_cart\".\"id\" = `).
			WithArgs(Cart["status"], test.AnyTime{}, Cart["user_id"], 1).
			WillReturnError(err)
		mock.ExpectRollback()
	} else {
		mock.ExpectExec(`UPDATE \"dp_cart\" SET (.+)  WHERE (.+) \"dp_cart\".\"id\" = `).
			WithArgs(Cart["status"], test.AnyTime{}, Cart["user_id"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`INSERT INTO "dp_cart_item"`).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
}
