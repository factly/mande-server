package cart

import (
	"database/sql/driver"
	"errors"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var CartItem map[string]interface{} = map[string]interface{}{
	"status":        "teststatus",
	"product_id":    1,
	"membership_id": 1,
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
		"status":        "teststatus1",
		"product_id":    1,
		"membership_id": 1,
	},
	{
		"status":        "teststatus2",
		"product_id":    1,
		"membership_id": 1,
	},
}

var CartItemCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "status", "user_id", "product_id", "membership_id"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_cart_item"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(1) FROM "dp_cart_item"`)
var errCartItemProductFK error = errors.New(`pq: insert or update on table "dp_cart_item" violates foreign key constraint "dp_cart_item_product_id_dp_product_id_foreign"`)

const basePath string = "/cartitems"
const path string = "/cartitems/{cartitem_id}"

func CartItemSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(selectQuery).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(CartItemCols).
			AddRow(1, time.Now(), time.Now(), nil, CartItem["status"], CartItem["user_id"], CartItem["product_id"], CartItem["membership_id"]))
}
