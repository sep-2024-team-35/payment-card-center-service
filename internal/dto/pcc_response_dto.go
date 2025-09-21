package dto

type PCCResponseDTO struct {
	Status            string `json:"status"`
	AcquirerOrderID   string `json:"acquirerOrderId"`
	AcquirerTimestamp string `json:"acquirerTimestamp"`
	IssuerOrderID     string `json:"issuerOrderId"`
	IssuerTimestamp   string `json:"issuerTimestamp"`
}
