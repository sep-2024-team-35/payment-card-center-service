package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/sep-2024-team-35/payment-card-center-service/internal/dto"
	"github.com/sep-2024-team-35/payment-card-center-service/internal/service"
)

type TransactionHandler struct {
	pcc *service.PCCService
}

func NewTransactionHandler(pcc *service.PCCService) *TransactionHandler {
	return &TransactionHandler{pcc: pcc}
}

// Execute godoc
// @Summary Submit interbank transaction request
// @Description Receives a payment request from Acquirer bank and routes it to the Issuer bank
// @Tags Transactions
// @Accept json
// @Produce json
// @Param request body dto.PaymentRequestDTO true "Payment request payload"
// @Success 200 {object} map[string]string "Transaction successfully routed"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Bank not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/transactions [post]
func (h *TransactionHandler) Execute(w http.ResponseWriter, r *http.Request) {
	log.Printf("⟳ Received /transactions request from %s", r.RemoteAddr)
	defer closeRequestBody(r)

	req, err := decodePaymentRequest(r)
	if err != nil {
		log.Printf("✗ Decode error: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	log.Printf("✔ Parsed request: OrderID=%s, PAN=%s",
		req.AcquirerOrderID, maskPAN(req.PrimaryAccountNumber),
	)

	bankID, err := extractBankID(req.PrimaryAccountNumber)
	if err != nil {
		log.Printf("✗ Invalid PAN: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid primaryAccountNumber")
		return
	}

	response, err := h.pcc.RouteToIssuer(bankID, req)
	if err != nil {
		log.Printf("✗ Routing failed: bankID=%s, error=%v", bankID, err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusOK, response); err != nil {
		log.Printf("✗ Encode response error: %v", err)
	} else {
		log.Printf("✓ Response sent: OrderID=%s, Status=%s", response.AcquirerOrderID, response.Status)
	}
}

func decodePaymentRequest(r *http.Request) (dto.ExternalTransactionRequestDTO, error) {
	var req dto.ExternalTransactionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}
	return req, nil
}

func extractBankID(pan string) (string, error) {
	if len(pan) < 16 {
		return "", errors.New("PAN too short")
	}
	return pan[:6], nil
}

func closeRequestBody(r *http.Request) {
	if err := r.Body.Close(); err != nil {
		log.Printf("failed to close request body: %v", err)
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func maskPAN(pan string) string {
	n := len(pan)
	if n <= 4 {
		return strings.Repeat("*", n)
	}
	return strings.Repeat("*", n-4) + pan[n-4:]
}
