package service

import (
	"fmt"
	"payment-card-center-service/internal/model"
	"github.com/sep-2024-team-35/payment-card-service/"
)

type PCCService struct {
	repo *repository.BankRepository
	// idempotencyKeys map[string]struct{} // kasnije za idempotentnost
}

func NewPCCService(repo *repository.BankRepository) *PCCService {
	return &PCCService{repo: repo}
}

// Rutiraj transakciju na Issuer bazirano na bank ID
func (s *PCCService) RouteToIssuer(bankID string, payload interface{}) (*model.Bank, error) {
	bank, err := s.repo.FindByID(bankID)
	if err != nil {
		return nil, fmt.Errorf("routing failed: %w", err)
	}
	// Ovdje bismo ubacili idempotency check po acquirerOrderID
	return bank, nil
}
