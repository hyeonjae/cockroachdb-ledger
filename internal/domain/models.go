package domain

import (
	"time"
)

type Account struct {
	ID            int       `json:"id" db:"id"`
	AccountNumber string    `json:"account_number" db:"account_number"`
	Balance       float64   `json:"balance" db:"balance"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type Holding struct {
	ID        int       `json:"id" db:"id"`
	AccountID int       `json:"account_id" db:"account_id"`
	StockCode string    `json:"stock_code" db:"stock_code"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Order struct {
	ID             int       `json:"id" db:"id"`
	AccountID      int       `json:"account_id" db:"account_id"`
	StockCode      string    `json:"stock_code" db:"stock_code"`
	Type           string    `json:"type" db:"type"`
	Direction      string    `json:"direction" db:"direction"`
	Quantity       int       `json:"quantity" db:"quantity"`
	Price          float64   `json:"price" db:"price"`
	FilledQuantity int       `json:"filled_quantity" db:"filled_quantity"`
	Status         string    `json:"status" db:"status"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CreateOrderRequest struct {
	AccountID int     `json:"account_id"`
	StockCode string  `json:"stock_code"`
	Type      string  `json:"type"`
	Direction string  `json:"direction"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type BalanceResponse struct {
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
}

type HoldingResponse struct {
	StockCode string `json:"stock_code"`
	Quantity  int    `json:"quantity"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}