package currency

import (
	"database/sql/driver"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var headers = map[string]string{
	"X-User":         "1",
	"X-Organisation": "1",
}

var Currency map[string]interface{} = map[string]interface{}{
	"name":     "Indian Rupees",
	"iso_code": "INR",
}

var undecodableCurrency map[string]interface{} = map[string]interface{}{
	"name":     1,
	"iso_code": 10,
}

var invalidCurrency map[string]interface{} = map[string]interface{}{
	"nam":     "Indian Rupee",
	"isocode": "Test ISO Code",
}

var CurrencyCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "created_by_id", "updated_by_id", "iso_code", "name"}

var selectQuery string = `SELECT (.+) FROM "dp_currency"`
var countQuery string = regexp.QuoteMeta(`SELECT count(1) FROM "dp_currency"`)

const basePath string = "/currencies"
const path string = "/currencies/{currency_id}"
const defaultPath string = "/currencies/default"

func CurrencySelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(selectQuery).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(CurrencyCols).
			AddRow(1, time.Now(), time.Now(), nil, 1, 1, Currency["iso_code"], Currency["name"]))
}

func currencyPaymentExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "dp_payment`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func currencyProductExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "dp_product`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func currencyDatasetExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "dp_dataset`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}
