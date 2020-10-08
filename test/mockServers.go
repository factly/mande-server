package test

import (
	"net/http"

	"github.com/factly/data-portal-server/util/razorpay"
	"github.com/spf13/viper"

	"gopkg.in/h2non/gock.v1"
)

var ReturnUpdate = map[string]interface{}{
	"updateId": 1,
}

var MeiliHits = map[string]interface{}{
	"hits": []map[string]interface{}{
		{
			"object_id":   "format_2",
			"kind":        "format",
			"id":          2,
			"description": "This is a test format",
			"name":        "Test Format",
			"is_default":  true,
		},
		{
			"object_id":   "format_3",
			"kind":        "format",
			"id":          3,
			"description": "This is second test format",
			"name":        "Test format 2",
			"is_default":  true,
		},
		{
			"object_id": "tag_2",
			"kind":      "tag",
			"id":        2,
			"slug":      "test-tag",
			"title":     "Test tag",
		},
	},
	"offset":           0,
	"limit":            10,
	"nbHits":           10,
	"exhaustiveNbHits": false,
	"processingTimeMs": 2,
	"query":            "test",
}

func MeiliGock() {
	gock.New(viper.GetString("meili.url")).
		Post("/indexes/data-portal/search").
		HeaderPresent("X-Meili-API-Key").
		Persist().
		Reply(http.StatusOK).
		JSON(MeiliHits)

	gock.New(viper.GetString("meili.url")).
		Post("/indexes/data-portal/documents").
		HeaderPresent("X-Meili-API-Key").
		Persist().
		Reply(http.StatusAccepted).
		JSON(ReturnUpdate)

	gock.New(viper.GetString("meili.url")).
		Put("/indexes/data-portal/documents").
		HeaderPresent("X-Meili-API-Key").
		Persist().
		Reply(http.StatusAccepted).
		JSON(ReturnUpdate)

	gock.New(viper.GetString("meili.url")).
		Delete("/indexes/data-portal/documents/(.+)").
		HeaderPresent("X-Meili-API-Key").
		Persist().
		Reply(http.StatusAccepted).
		JSON(ReturnUpdate)
}

var RazorpayOrder = map[string]interface{}{
	"id":          "order_FltCdu23fGaTwG",
	"entity":      "order",
	"amount":      5000,
	"amount_paid": 0,
	"amount_due":  5000,
	"currency":    "INR",
	"receipt":     "Test Receipt no. 1",
	"offer_id":    nil,
	"status":      "created",
	"attempts":    0,
	"notes": map[string]interface{}{
		"info": "this payment is for first order",
	},
	"created_at": 1602047090,
}

var RazorpayPayment = map[string]interface{}{
	"id":              "pay_FjYWQFwuiE89Xp",
	"entity":          "payment",
	"amount":          10000,
	"currency":        "INR",
	"status":          "captured",
	"order_id":        "order_FjYVOJ8Vod4lmT",
	"invoice_id":      nil,
	"international":   false,
	"method":          "card",
	"amount_refunded": 0,
	"refund_status":   nil,
	"captured":        true,
	"description":     "Test Transaction",
	"card_id":         "card_FjYNqO7cTrB4EU",
	"bank":            nil,
	"wallet":          nil,
	"vpa":             nil,
	"email":           "gaurav.kumar@example.com",
	"contact":         "+919999999999",
	"notes": map[string]interface{}{
		"address": "Razorpay Corporate Office",
	},
	"fee":               2798,
	"tax":               0,
	"error_code":        nil,
	"error_description": nil,
	"error_source":      nil,
	"error_step":        nil,
	"error_reason":      nil,
	"acquirer_data": map[string]interface{}{
		"auth_code": "464641",
	},
	"created_at": 1601889873,
}

func RazorpayGock() {
	razorpay.SetupClient()
	viper.Set("razorpay.secret", "testsecret")

	gock.New("https://api.razorpay.com").
		Post("/v1/orders").
		Persist().
		Reply(http.StatusOK).
		JSON(RazorpayOrder)

	gock.New("https://api.razorpay.com").
		Get("/v1/payments/(.+)").
		Persist().
		Reply(http.StatusOK).
		JSON(RazorpayPayment)
}
