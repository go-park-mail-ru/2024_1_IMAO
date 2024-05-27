package models

import "time"

// Payment "pending"
type PaymentList struct {
	Type  string                     `json:"type"`
	Items []WaitingForCapturePayment `json:"items"`
}

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
	Duration    string `json:"duration"`
}

type ReceivedPaymentFormItem struct {
	AdvertId uint `json:"advertId"`
	Rate     uint `json:"rate"`
}

type PaymnetUuidList struct {
	UuidList []string `json:"uuidList"`
}

type PaymnetUuidListPad struct {
	Pad []*string `json:"pad"`
}

type PaymentsDatesList struct {
	List []*time.Time
}

// Payment "waiting_for_capture"
type WaitingForCapturePayment struct {
	ID                   string                                `json:"id"`
	Status               string                                `json:"status"`
	Amount               WaitingForCapturePaymentAmount        `json:"amount"`
	Description          string                                `json:"description"`
	Recipient            WaitingForCapturePaymentRecipient     `json:"recipient"`
	PaymentMethod        WaitingForCapturePaymentPaymentMethod `json:"payment_method"`
	CreatedAt            string                                `json:"created_at"`
	ExpiresAt            string                                `json:"expires_at"`
	Test                 bool                                  `json:"test"`
	Paid                 bool                                  `json:"paid"`
	Refundable           bool                                  `json:"refundable"`
	Metadata             map[string]interface{}                `json:"metadata"`
	AuthorizationDetails AuthorizationDetails                  `json:"authorization_details"`
}

type WaitingForCapturePaymentAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type WaitingForCapturePaymentRecipient struct {
	AccountID string `json:"account_id"`
	GatewayID string `json:"gateway_id"`
}

type WaitingForCapturePaymentPaymentMethod struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Saved bool   `json:"saved"`
	Title string `json:"title"`
	Card  Card   `json:"card"`
}

type Card struct {
	First6        string      `json:"first6"`
	Last4         string      `json:"last4"`
	ExpiryYear    string      `json:"expiry_year"`
	ExpiryMonth   string      `json:"expiry_month"`
	CardType      string      `json:"card_type"`
	CardProduct   CardProduct `json:"card_product"`
	IssuerCountry string      `json:"issuer_country"`
}

type CardProduct struct {
	Code string `json:"code"`
}

type AuthorizationDetails struct {
	RRN          string       `json:"rrn"`
	AuthCode     string       `json:"auth_code"`
	ThreeDSecure ThreeDSecure `json:"three_d_secure"`
}

type ThreeDSecure struct {
	Applied            bool   `json:"applied"`
	Protocol           string `json:"protocol"`
	MethodCompleted    bool   `json:"method_completed"`
	ChallengeCompleted bool   `json:"challenge_completed"`
}

type PaymentFormResponse struct {
	PaymentFormUrl string `json:"paymentFormUrl"`
}
