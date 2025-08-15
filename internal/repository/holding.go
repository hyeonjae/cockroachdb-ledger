package repository

import (
	"database/sql"
	"mini-ledger/internal/db"
	"mini-ledger/internal/domain"
)

type holdingRepository struct{}

func NewHoldingRepository() HoldingRepository {
	return &holdingRepository{}
}

func (r *holdingRepository) GetByAccountID(querier db.Querier, accountID int) ([]*domain.Holding, error) {
	var holdings []*domain.Holding
	query := `SELECT id, account_id, stock_code, quantity, created_at, updated_at FROM holdings WHERE account_id = $1`
	err := querier.Select(&holdings, query, accountID)
	if err != nil {
		return nil, err
	}
	return holdings, nil
}

func (r *holdingRepository) GetByAccountIDAndStockCode(querier db.Querier, accountID int, stockCode string) (*domain.Holding, error) {
	var holding domain.Holding
	query := `SELECT id, account_id, stock_code, quantity, created_at, updated_at FROM holdings WHERE account_id = $1 AND stock_code = $2`
	err := querier.Get(&holding, query, accountID, stockCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &holding, nil
}

func (r *holdingRepository) UpdateQuantity(querier db.Querier, accountID int, stockCode string, quantity int) error {
	if quantity <= 0 {
		query := `DELETE FROM holdings WHERE account_id = $1 AND stock_code = $2`
		_, err := querier.Exec(query, accountID, stockCode)
		return err
	}

	query := `UPDATE holdings SET quantity = $1, updated_at = NOW() WHERE account_id = $2 AND stock_code = $3`
	_, err := querier.Exec(query, quantity, accountID, stockCode)
	return err
}

func (r *holdingRepository) Create(querier db.Querier, holding *domain.Holding) error {
	query := `INSERT INTO holdings (account_id, stock_code, quantity) VALUES ($1, $2, $3) 
			  ON CONFLICT (account_id, stock_code) DO UPDATE SET quantity = holdings.quantity + EXCLUDED.quantity`
	_, err := querier.Exec(query, holding.AccountID, holding.StockCode, holding.Quantity)
	return err
}