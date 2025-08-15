package domain

import "errors"

var (
	ErrAccountNotFound           = errors.New("account not found")
	ErrOrderNotFound             = errors.New("order not found")
	ErrInsufficientFunds         = errors.New("insufficient funds")
	ErrInsufficientHoldingQuantity = errors.New("insufficient holding quantity")
	ErrOrderNotCancelable        = errors.New("order is not in a cancelable state")
)