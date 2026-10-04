package service

import (
	"fmt"

	"github.com/shashank/order-service/internal/dto"
	apperrors "github.com/shashank/order-service/internal/errors"
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

		valiationError := &apperrors.ValidationErrors{
			Field:   "amount",
			Message: "Amount is invalid",
		}
		return nil, fmt.Errorf("order validation failed : %w", valiationError)
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
	ord, err := o.orderRepository.GetOrderByID(id)
	if err != nil {
		return nil, fmt.Errorf("Order not found with id %s : %w", id, err)
	}
	return ord, nil
}

func (o *OrderService) DeleteOrderByID(id string) (string, error) {

	message, err := o.orderRepository.DeleteOrderByID(id)
	if err != nil {
		return "", fmt.Errorf("Order not found with order id %s :%w", id, err)
	}
	return message, nil
}

func (o *OrderService) UpdateOrderByID(id string, request *dto.UpdateOrderRequest) (*model.Order, error) {

	order := model.Order{
		Status: request.Status,
	}

	existingOrders, err := o.orderRepository.GetOrderByID(id)

	if err != nil {
		return nil, fmt.Errorf("Order not found with order id %s : %w", id, err)
	}

	if !isValidTransition(existingOrders.Status, request.Status) {
		return nil, fmt.Errorf("Invalid status")
	}

	updatedOrder, err := o.orderRepository.UpdateOrderByID(id, &order)

	if err != nil {
		return nil, fmt.Errorf("Order not found with id %s : %w", id, err)
	}

	return updatedOrder, nil
}

func isValidTransition(existingStatus string, newStatus string) bool {

	if existingStatus == "CREATED" {
		if newStatus == "CONFIRMED" {
			return true
		}
	}

	if existingStatus == "CONFIRMED" {
		if newStatus == "SHIPPED" {
			return true
		}
	}

	if existingStatus == "SHIPPED" {
		if newStatus == "DELIVERED" {
			return true
		}
	}

	return false
}
