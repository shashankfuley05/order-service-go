package handler

import (
	"github.com/shashank/order-service/internal/dto"
	"github.com/shashank/order-service/internal/model"
)

type OrderService interface {
	CreateOrder(request *dto.CreateOrderRequest) (*model.Order, error)
	FetchOrders() []model.Order
	GetOrderByID(id string) (*model.Order, error)
	DeleteOrderByID(id string) (string, error)
	UpdateOrderByID(id string, request *dto.UpdateOrderRequest) (*model.Order, error)
}
