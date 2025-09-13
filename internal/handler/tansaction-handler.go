package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"payment-card-center-service/internal/service"
)

type TransactionHandler struct {
	pcc *service.PCCService
}

func NewTransactionHandler(pcc *service.PCCService) *TransactionHandler {
	return &TransactionHandler{pcc: pcc}
}

func (h *TransactionHandler) Execute(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var req struct {
		AcquirerOrderID string  `json:"acquirerOrderId"`
		AcquirerTS      string  `json:"acquirerTimestamp"`
		Pan             string  `json:"pan"`
		Amount          float64 `json:"amount"`
		Currency        string  `json:"currency"`
		// TransactionDto itd...
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Ekstraktujemo bankID iz PAN-a (npr. prvih 4 cifre)
	bankID := req.Pan[:4]
	_, err := h.pcc.RouteToIssuer(bankID, req)
	if err != nil {
		http.Error(w, "bank not found", http.StatusNotFound)
		return
	}

	// TODO: HTTP client ka bank.URL sa timeout i retry mehanizmom

	resp := struct {
		AcquirerOrderID string `json:"acquirerOrderId"`
		Status          string `json:"status"`
		// issuerOrderId, timestamps...
	}{
		AcquirerOrderID: req.AcquirerOrderID,
		Status:          "Pending",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
