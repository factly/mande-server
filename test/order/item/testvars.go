package orderitem

import (
	"errors"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var OrderItem map[string]interface{} = map[string]interface{}{
	"extra_info": "Test Extra Info",
	"product_id": 1,
}

var invalidOrderItem map[string]interface{} = map[string]interface{}{
	"extra_info": "Test Extra Info",
	"productid":  1,
}

var undecodableOrderItem map[string]interface{} = map[string]interface{}{
	"extra_info": 23,
	"productid":  "1",
}

var orderitemlist []map[string]interface{} = []map[string]interface{}{
	{
		"extra_info": "Test Extra Info 1",
		"product_id": 1,
	},
	{
		"extra_info": "Test Extra Info 2",
		"product_id": 1,
	},
}

var OrderItemCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "extra_info", "product_id", "order_id"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_order_item"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_order_item"`)
var errOrderItemFK error = errors.New(`pq: insert or update on table "dp_order_item" violates foreign key constraint "dp_order_item_product_id_dp_product_id_foreign"`)

const basePath string = "/orders/{order_id}/items"
const path string = "/orders/{order_id}/items/{item_id}"

func OrderItemSelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(OrderItemCols).
			AddRow(1, time.Now(), time.Now(), nil, OrderItem["extra_info"], OrderItem["product_id"], 1))
}
