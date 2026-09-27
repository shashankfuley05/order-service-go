package apperrors

import "errors"

var (
	ErrOrderNotFound  = errors.New("Order not found")
	ErrDuplicateOrder = errors.New("Order already exist")
	ErrInvalidAmount  = errors.New("Amount is invalid")
)

type ValidationErrors struct {
	Field   string
	Message string
}

func (err *ValidationErrors) Error() string {
	return err.Message
}
