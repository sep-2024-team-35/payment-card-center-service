package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/sep-2024-team-35/payment-card-center-service/internal/dto"
	"github.com/sep-2024-team-35/payment-card-center-service/internal/repository"
)

type PCCService struct {
	repo             *repository.BankRepository
	mu               sync.RWMutex
	idempotencyStore map[string]*dto.PCCResponseDTO
}

// Inicijalizacija servisa sa praznim idempotentnim skladištem
func NewPCCService(repo *repository.BankRepository) *PCCService {
	return &PCCService{
		repo:             repo,
		idempotencyStore: make(map[string]*dto.PCCResponseDTO),
	}
}

// RouteToIssuer sada prima ceo DTO da bi mogao koristiti AcquirerOrderID
func (s *PCCService) RouteToIssuer(bankID string, req dto.ExternalTransactionRequestDTO) (*dto.PCCResponseDTO, error) {
	// 1. Idempotency check
	s.mu.RLock()
	if previous, found := s.idempotencyStore[req.AcquirerOrderID]; found {
		s.mu.RUnlock()
		log.Printf("[INFO] Idempotent request detected: OrderID=%s", req.AcquirerOrderID)
		return previous, nil
	}
	s.mu.RUnlock()

	// 2. Nađi banku
	bank, err := s.repo.FindByID(bankID)
	if err != nil {
		log.Printf("[ERROR] Bank not found for ID=%s", bankID)
		return nil, fmt.Errorf("routing failed: %w", err)
	}

	// 3. Marshal payload
	payload, err := json.Marshal(req)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal request: %v", err)
		return nil, fmt.Errorf("invalid request payload: %w", err)
	}

	// 4. Poziv ka Issuer banci
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(bank.URL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("[ERROR] Failed to contact Issuer bank: %v", err)
		return nil, fmt.Errorf("issuer bank unreachable: %w", err)
	}
	if resp == nil {
		log.Printf("[ERROR] Issuer bank returned nil response")
		return nil, fmt.Errorf("issuer bank returned no response")
	}
	if resp.Body != nil {
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				log.Printf("[WARN] Failed to close response body: %v", cerr)
			}
		}()
	}

	// 5. Proveri status kod
	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Issuer bank responded with status %d", resp.StatusCode)
		return nil, fmt.Errorf("issuer bank returned status %d", resp.StatusCode)
	}

	// 6. Decode response
	var issuerResponse dto.PCCResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&issuerResponse); err != nil {
		log.Printf("[ERROR] Failed to decode issuer response: %v", err)
		return nil, fmt.Errorf("invalid issuer response: %w", err)
	}

	// 7. Čuvanje u idempotencyStore
	s.mu.Lock()
	s.idempotencyStore[req.AcquirerOrderID] = &issuerResponse
	s.mu.Unlock()

	log.Printf("[INFO] Issuer response stored for OrderID=%s", req.AcquirerOrderID)
	return &issuerResponse, nil
}
