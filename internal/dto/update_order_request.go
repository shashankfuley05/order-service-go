package dto

type UpdateOrderRequest struct {
	CustomerName string `json:"customerName" binding:"required"`

	Amount float64 `json:"amount" binding:"required,gt=0"`
}
