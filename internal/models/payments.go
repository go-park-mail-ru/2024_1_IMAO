package models

// Payment "pending"
type Payment struct {
	ID            string                 `json:"id"`
	Status        string                 `json:"status"`
	Amount        Amount                 `json:"amount"`
	Description   string                 `json:"description"`
	Recipient     Recipient              `json:"recipient"`
	PaymentMethod PaymentMethod          `json:"payment_method"`
	CreatedAt     string                 `json:"created_at"`
	Confirmation  Confirmation           `json:"confirmation"`
	Test          bool                   `json:"test"`
	Paid          bool                   `json:"paid"`
	Refundable    bool                   `json:"refundable"`
	Metadata      map[string]interface{} `json:"metadata"`
}

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Recipient struct {
	AccountID string `json:"account_id"`
	GatewayID string `json:"gateway_id"`
}

type PaymentMethod struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Saved bool   `json:"saved"`
}

type Confirmation struct {
	Type            string `json:"type"`
	ReturnURL       string `json:"return_url"`
	ConfirmationURL string `json:"confirmation_url"`
}

// PaymentInitData
type PaymentInitData struct {
	Amount            PaymentInitAmount            `json:"amount"`
	PaymentMethodData PaymentInitPaymentMethodData `json:"payment_method_data"`
	Confirmation      PaymentInitConfirmation      `json:"confirmation"`
	Description       string                       `json:"description"`
}

type PaymentInitAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type PaymentInitPaymentMethodData struct {
	Type string `json:"type"`
}

type PaymentInitConfirmation struct {
	Type      string `json:"type"`
	ReturnURL string `json:"return_url"`
}

// Utils structures

type PriceAndDescription struct {
	Price       string `json:"price"`
	Description string `json:"description"`
	UrlEnding   string `json:"urlEnding"`
}

type ReceivedPaymentFormItem struct {
	AdvertId uint `json:"advertId"`
	Rate     uint `json:"rate"`
}
