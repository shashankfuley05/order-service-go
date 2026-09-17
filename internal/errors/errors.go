package errors

import "errors"

var (
	ErrOrderNotFound  = errors.New("Order not found")
	ErrDuplicateOrder = errors.New("Order already exist")
	ErrInvalidAmount  = errors.New("Amount is invalid")
)
