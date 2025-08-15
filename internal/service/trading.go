package service

import (
	"database/sql"
	"mini-ledger/internal/db"
	"mini-ledger/internal/domain"
	"mini-ledger/internal/repository"
)

type TradingService struct {
	db              *db.Database
	accountRepo     repository.AccountRepository
	holdingRepo     repository.HoldingRepository
	orderRepo       repository.OrderRepository
}

func NewTradingService(
	database *db.Database,
	accountRepo repository.AccountRepository,
	holdingRepo repository.HoldingRepository,
	orderRepo repository.OrderRepository,
) *TradingService {
	return &TradingService{
		db:          database,
		accountRepo: accountRepo,
		holdingRepo: holdingRepo,
		orderRepo:   orderRepo,
	}
}

func (s *TradingService) GetAccountBalance(accountID int) (*domain.BalanceResponse, error) {
	account, err := s.accountRepo.GetByID(s.db, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAccountNotFound
		}
		return nil, err
	}

	return &domain.BalanceResponse{
		AccountNumber: account.AccountNumber,
		Balance:       account.Balance,
	}, nil
}

func (s *TradingService) GetAccountHoldings(accountID int) ([]*domain.HoldingResponse, error) {
	_, err := s.accountRepo.GetByID(s.db, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAccountNotFound
		}
		return nil, err
	}

	holdings, err := s.holdingRepo.GetByAccountID(s.db, accountID)
	if err != nil {
		return nil, err
	}

	var response []*domain.HoldingResponse
	for _, holding := range holdings {
		response = append(response, &domain.HoldingResponse{
			StockCode: holding.StockCode,
			Quantity:  holding.Quantity,
		})
	}

	return response, nil
}

func (s *TradingService) CreateOrder(req *domain.CreateOrderRequest) (*domain.Order, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	account, err := s.accountRepo.GetByID(tx, req.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAccountNotFound
		}
		return nil, err
	}

	if req.Direction == "BUY" {
		totalCost := req.Price * float64(req.Quantity)
		if account.Balance < totalCost {
			return nil, domain.ErrInsufficientFunds
		}

		newBalance := account.Balance - totalCost
		if err := s.accountRepo.UpdateBalance(tx, req.AccountID, newBalance); err != nil {
			return nil, err
		}
	} else if req.Direction == "SELL" {
		holding, err := s.holdingRepo.GetByAccountIDAndStockCode(tx, req.AccountID, req.StockCode)
		if err != nil {
			return nil, err
		}
		if holding == nil || holding.Quantity < req.Quantity {
			return nil, domain.ErrInsufficientHoldingQuantity
		}

		newQuantity := holding.Quantity - req.Quantity
		if err := s.holdingRepo.UpdateQuantity(tx, req.AccountID, req.StockCode, newQuantity); err != nil {
			return nil, err
		}
	}

	order := &domain.Order{
		AccountID:      req.AccountID,
		StockCode:      req.StockCode,
		Type:           req.Type,
		Direction:      req.Direction,
		Quantity:       req.Quantity,
		Price:          req.Price,
		FilledQuantity: 0,
		Status:         "PENDING",
	}

	createdOrder, err := s.orderRepo.Create(tx, order)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return createdOrder, nil
}

func (s *TradingService) CancelOrder(orderID int) (*domain.Order, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	order, err := s.orderRepo.GetByID(tx, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	if order.Status != "PENDING" && order.Status != "PARTIAL" {
		return nil, domain.ErrOrderNotCancelable
	}

	unfilledQuantity := order.Quantity - order.FilledQuantity

	if order.Direction == "BUY" {
		refundAmount := order.Price * float64(unfilledQuantity)
		account, err := s.accountRepo.GetByID(tx, order.AccountID)
		if err != nil {
			return nil, err
		}
		newBalance := account.Balance + refundAmount
		if err := s.accountRepo.UpdateBalance(tx, order.AccountID, newBalance); err != nil {
			return nil, err
		}
	} else if order.Direction == "SELL" {
		holding, err := s.holdingRepo.GetByAccountIDAndStockCode(tx, order.AccountID, order.StockCode)
		if err != nil {
			return nil, err
		}
		
		var newQuantity int
		if holding == nil {
			newQuantity = unfilledQuantity
			newHolding := &domain.Holding{
				AccountID: order.AccountID,
				StockCode: order.StockCode,
				Quantity:  newQuantity,
			}
			if err := s.holdingRepo.Create(tx, newHolding); err != nil {
				return nil, err
			}
		} else {
			newQuantity = holding.Quantity + unfilledQuantity
			if err := s.holdingRepo.UpdateQuantity(tx, order.AccountID, order.StockCode, newQuantity); err != nil {
				return nil, err
			}
		}
	}

	if err := s.orderRepo.UpdateStatus(tx, orderID, "CANCELED"); err != nil {
		return nil, err
	}

	updatedOrder, err := s.orderRepo.GetByID(tx, orderID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return updatedOrder, nil
}