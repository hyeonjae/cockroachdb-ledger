package repository

import (
	"mini-ledger/internal/db"
	"mini-ledger/internal/domain"
)

type orderRepository struct{}

func NewOrderRepository() OrderRepository {
	return &orderRepository{}
}

func (r *orderRepository) Create(querier db.Querier, order *domain.Order) (*domain.Order, error) {
	query := `INSERT INTO orders (account_id, stock_code, type, direction, quantity, price, filled_quantity, status) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	
	var id int
	err := querier.Get(&id, query, order.AccountID, order.StockCode, order.Type, order.Direction, 
		order.Quantity, order.Price, order.FilledQuantity, order.Status)
	if err != nil {
		return nil, err
	}
	
	return r.GetByID(querier, id)
}

func (r *orderRepository) GetByID(querier db.Querier, id int) (*domain.Order, error) {
	var order domain.Order
	query := `SELECT id, account_id, stock_code, type, direction, quantity, price, filled_quantity, status, created_at, updated_at 
			  FROM orders WHERE id = $1`
	err := querier.Get(&order, query, id)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateStatus(querier db.Querier, id int, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := querier.Exec(query, status, id)
	return err
}