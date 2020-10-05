package razorpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/factly/data-portal-server/config"
)

// VerifySignature verifies razorpay payment signature
func VerifySignature(orderID, paymentID, signature string) bool {
	h := hmac.New(sha256.New, []byte(config.RazorpaySecret))
	h.Write([]byte(fmt.Sprint(orderID, "|", paymentID)))

	generated := hex.EncodeToString(h.Sum(nil))

	return generated == signature
}
