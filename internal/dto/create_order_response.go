package dto

type CreateOrderResponse struct {
	CustomerName string  `json:"customerName"`
	Amount       float64 `json:"amount"`
	ID           string  `json:"id"`
	Status       string  `json:"status"`
}
