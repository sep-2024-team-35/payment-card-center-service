package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ExternalTransactionRequestDTO struct {
	// Transaction Data
	ID                uuid.UUID       `json:"id"`
	AcquirerOrderID   string          `json:"acquirerOrderId"`
	AcquirerTimestamp string          `json:"acquirerTimestamp"`
	Amount            decimal.Decimal `json:"amount"`
	MerchantOrderID   string          `json:"merchantOrderId"`
	MerchantTimestamp string          `json:"merchantTimestamp"`
	Currency          string          `json:"currency"`

	// Card Data
	PrimaryAccountNumber string    `json:"primaryAccountNumber"`
	CardHolderName       string    `json:"cardHolderName,omitempty"`
	ExpirationDate       string    `json:"expirationDate,omitempty"`
	SecurityCode         string    `json:"securityCode,omitempty"`
	PaymentRequestID     uuid.UUID `json:"paymentRequestId,omitempty"`
}
