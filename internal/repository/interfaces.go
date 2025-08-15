package repository

import (
	"mini-ledger/internal/db"
	"mini-ledger/internal/domain"
)

type AccountRepository interface {
	GetByID(querier db.Querier, id int) (*domain.Account, error)
	UpdateBalance(querier db.Querier, id int, balance float64) error
}

type HoldingRepository interface {
	GetByAccountID(querier db.Querier, accountID int) ([]*domain.Holding, error)
	GetByAccountIDAndStockCode(querier db.Querier, accountID int, stockCode string) (*domain.Holding, error)
	UpdateQuantity(querier db.Querier, accountID int, stockCode string, quantity int) error
	Create(querier db.Querier, holding *domain.Holding) error
}

type OrderRepository interface {
	Create(querier db.Querier, order *domain.Order) (*domain.Order, error)
	GetByID(querier db.Querier, id int) (*domain.Order, error)
	UpdateStatus(querier db.Querier, id int, status string) error
}