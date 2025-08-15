package repository

import (
	"mini-ledger/internal/db"
	"mini-ledger/internal/domain"
)

type accountRepository struct{}

func NewAccountRepository() AccountRepository {
	return &accountRepository{}
}

func (r *accountRepository) GetByID(querier db.Querier, id int) (*domain.Account, error) {
	var account domain.Account
	query := `SELECT id, account_number, balance, created_at, updated_at FROM accounts WHERE id = $1`
	err := querier.Get(&account, query, id)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) UpdateBalance(querier db.Querier, id int, balance float64) error {
	query := `UPDATE accounts SET balance = $1, updated_at = NOW() WHERE id = $2`
	_, err := querier.Exec(query, balance, id)
	return err
}