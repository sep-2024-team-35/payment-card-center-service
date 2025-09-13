package repository

import (
	"errors"
	"github.com/sep-2024-team-35/payment-card-center-service/internal/config"
	"github.com/sep-2024-team-35/payment-card-center-service/internal/model"
)

type BankRepository struct {
	banks map[string]*model.Bank
}

func NewBankRepository() *BankRepository {
	repo := &BankRepository{banks: make(map[string]*model.Bank)}
	for _, b := range config.Global.Banks {
		repo.banks[b.ID] = &model.Bank{ID: b.ID, Name: b.Name, URL: b.URL}
	}
	return repo
}

// Pronadji banku po tačnom ID (npr. 4-znamenkasti BIN)
func (r *BankRepository) FindByID(id string) (*model.Bank, error) {
	if bank, ok := r.banks[id]; ok {
		return bank, nil
	}
	return nil, errors.New("bank not found")
}
