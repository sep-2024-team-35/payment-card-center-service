package service

import (
	"fmt"
	"sync"

	"github.com/sep-2024-team-35/payment-card-center-service/internal/dto"
	"github.com/sep-2024-team-35/payment-card-center-service/internal/model"
	"github.com/sep-2024-team-35/payment-card-center-service/internal/repository"
)

type PCCService struct {
	repo             *repository.BankRepository
	mu               sync.RWMutex
	idempotencyStore map[string]*model.Bank
}

// Inicijalizacija servisa sa praznim idempotentnim skladištem
func NewPCCService(repo *repository.BankRepository) *PCCService {
	return &PCCService{
		repo:             repo,
		idempotencyStore: make(map[string]*model.Bank),
	}
}

// RouteToIssuer sada prima ceo DTO da bi mogao koristiti AcquirerOrderID
func (s *PCCService) RouteToIssuer(bankID string, req dto.PaymentRequestDTO) (*model.Bank, error) {
	// 1. IDempotency check
	s.mu.RLock()
	if previous, found := s.idempotencyStore[req.AcquirerOrderID]; found {
		s.mu.RUnlock()
		return previous, nil
	}
	s.mu.RUnlock()

	// 2. Osnovno rutiranje na banku iz repozitorijuma
	bank, err := s.repo.FindByID(bankID)
	if err != nil {
		return nil, fmt.Errorf("routing failed: %w", err)
	}

	// 3. (ovde će ići HTTP poziv i obrada odgovora Issuer-a)

	// 4. Čuvanje u idempotentno skladište
	s.mu.Lock()
	s.idempotencyStore[req.AcquirerOrderID] = bank
	s.mu.Unlock()

	return bank, nil
}
