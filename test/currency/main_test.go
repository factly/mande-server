package currency

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var currency map[string]interface{} = map[string]interface{}{
	"name":     "Test Name",
	"iso_code": "Test ISO Code",
}

var invalidCurrency map[string]interface{} = map[string]interface{}{
	"nam":     "Test Name",
	"isocode": "Test ISO Code",
}

var currencyCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "iso_code", "name"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_currency"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_currency"`)

const basePath string = "/currencies"
const path string = "/currencies/{currency_id}"

func currencySelectMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(currencyCols).
			AddRow(1, time.Now(), time.Now(), nil, currency["iso_code"], currency["name"]))
}

func currencyPaymentExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_payment`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func currencyProductExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_product`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func currencyDatasetExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_dataset`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}
