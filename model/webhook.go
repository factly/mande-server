package model

// RazorpayPayment razorpay object
type RazorpayPayment struct {
	ID               string        `json:"id,omitempty"`
	Entity           string        `json:"entity,omitempty"`
	Amount           int           `json:"amount,omitempty"`
	Currency         string        `json:"currency,omitempty"`
	Status           string        `json:"status,omitempty"`
	OrderID          string        `json:"order_id,omitempty"`
	InvoiceID        interface{}   `json:"invoice_id,omitempty"`
	International    bool          `json:"international,omitempty"`
	Method           string        `json:"method,omitempty"`
	AmountRefunded   int           `json:"amount_refunded,omitempty"`
	RefundStatus     interface{}   `json:"refund_status,omitempty"`
	Captured         bool          `json:"captured,omitempty"`
	Description      interface{}   `json:"description,omitempty"`
	CardID           interface{}   `json:"card_id,omitempty"`
	Bank             string        `json:"bank,omitempty"`
	Wallet           interface{}   `json:"wallet,omitempty"`
	VPA              interface{}   `json:"vpa,omitempty"`
	Email            string        `json:"email,omitempty"`
	Contact          string        `json:"contact,omitempty"`
	Notes            []interface{} `json:"notes,omitempty"`
	Fee              interface{}   `json:"fee,omitempty"`
	Tax              interface{}   `json:"tax,omitempty"`
	ErrorCode        string        `json:"error_code,omitempty"`
	ErrorDescription string        `json:"error_description,omitempty"`
	ErrorSource      string        `json:"error_source,omitempty"`
	ErrorStep        string        `json:"error_step,omitempty"`
	ErrorReason      string        `json:"error_reason,omitempty"`
	AquirerData      interface{}   `json:"aquirer_data,omitempty"`
	CreatedAt        uint32        `json:"created_at,omitempty"`
}

// PaymentEntity razorpay object
type PaymentEntity struct {
	Entity RazorpayPayment `json:"entity,omitempty"`
}

// PaymentPayload razorpay object
type PaymentPayload struct {
	Payment PaymentEntity `json:"payment,omitempty"`
}

// FailedPaymentEvent razorpay object
type FailedPaymentEvent struct {
	Entity    string         `json:"entity,omitempty"`
	AccountID string         `json:"account_id,omitempty"`
	Event     string         `json:"event,omitempty"`
	Contains  []string       `json:"contains,omitempty"`
	Payload   PaymentPayload `json:"payload,omitempty"`
	CreatedAt uint32         `json:"created_at,omitempty"`
}
