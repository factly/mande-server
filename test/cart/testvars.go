package cart

import (
	"errors"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/test"
)

var CartItem map[string]interface{} = map[string]interface{}{
	"status":     "teststatus",
	"product_id": 1,
}

var invalidCartItem map[string]interface{} = map[string]interface{}{
	"stat":      "testtus",
	"productds": 1,
}

var undecodableCartItem map[string]interface{} = map[string]interface{}{
	"status":     45,
	"product_id": "1",
}

var cartitemslist []map[string]interface{} = []map[string]interface{}{
	{
		"status":     "teststatus1",
		"product_id": 1,
	},
	{
		"status":     "teststatus2",
		"product_id": 1,
	},
}

var CartItemCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "status", "user_id", "product_id"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_cart_item"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_cart_item"`)
var errCartItemProductFK error = errors.New(`pq: insert or update on table "dp_cart_item" violates foreign key constraint "dp_cart_item_product_id_dp_product_id_foreign"`)

const basePath string = "/cartitems"
const path string = "/cartitems/{cartitem_id}"

func CartItemSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(CartItemCols).
			AddRow(1, time.Now(), time.Now(), nil, CartItem["status"], CartItem["user_id"], CartItem["product_id"]))
}

func preUpdateMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(CartItemCols).
			AddRow(1, time.Now(), time.Now(), nil, "status", 1, 2))

	mock.ExpectBegin()
}

func updateMock(mock sqlmock.Sqlmock, err error) {
	preUpdateMock(mock)

	if err != nil {
		mock.ExpectExec(`UPDATE \"dp_cart_item\" SET (.+)  WHERE (.+) \"dp_cart_item\".\"id\" = `).
			WithArgs(CartItem["product_id"], CartItem["status"], test.AnyTime{}, 1, 1).
			WillReturnError(err)
		mock.ExpectRollback()
	} else {
		mock.ExpectExec(`UPDATE \"dp_cart_item\" SET (.+)  WHERE (.+) \"dp_cart_item\".\"id\" = `).
			WithArgs(CartItem["product_id"], CartItem["status"], test.AnyTime{}, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
