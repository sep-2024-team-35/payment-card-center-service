package dto

type PaymentRequestDTO struct {
	AcquirerOrderID      string  `json:"acquirerOrderId"`
	AcquirerTimestamp    string  `json:"acquirerTimestamp"`
	PrimaryAccountNumber string  `json:"primaryAccountNumber"`
	Amount               float64 `json:"amount"`
	Currency             string  `json:"currency"`
}
