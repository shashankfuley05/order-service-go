package dto

type UpdateOrderRequest struct {
	Status string `json:"status" binding:"required"`
}
