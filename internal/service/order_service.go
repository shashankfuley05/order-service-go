package service

import (
	"errors"

	"github.com/shashank/order-service/internal/dto"
	"github.com/shashank/order-service/internal/model"
	"github.com/shashank/order-service/internal/repository"
)

type OrderService struct {
	orderRepository repository.OrderRepository
}

func NewOrderService(r repository.OrderRepository) *OrderService {
	return &OrderService{
		orderRepository: r,
	}
}

func (o *OrderService) CreateOrder(request *dto.CreateOrderRequest) (*model.Order, error) {

	if request.Amount < 100000000.00 {
		return nil, errors.New("BC amount ye kya chillar amount hai")
	}
	order := &model.Order{
		CustomerName: request.CustomerName,
		Amount:       request.Amount,
	}

	return o.orderRepository.Save(order)

}

func (o *OrderService) FetchOrders() []model.Order {

	return o.orderRepository.GetAll()
}

func (o *OrderService) GetOrderByID(id string) (*model.Order, error) {
	return o.orderRepository.GetOrderByID(id)
}

func (o *OrderService) DeleteOrderByID(id string) (string, error) {
	return o.orderRepository.DeleteOrderByID(id)
}

func (o *OrderService) UpdateOrderByID(id string, request *dto.UpdateOrderRequest) (*model.Order, error) {

	order := model.Order{
		CustomerName: request.CustomerName,
		Amount:       request.Amount,
	}

	return o.orderRepository.UpdateOrderByID(id, &order)
}
