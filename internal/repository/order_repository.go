package repository

import (
	"fmt"
	"strings"

	"github.com/shashank/order-service/internal/model"
)

type InMemoryOrderRepository struct {
	orders  []model.Order
	counter int64
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		counter: 1,
	}
}

func (r *InMemoryOrderRepository) Save(o *model.Order) (*model.Order, error) {
	o.ID = fmt.Sprintf("ORD-%d", r.counter)
	o.Status = "CREATED"
	r.orders = append(r.orders, *o)

	r.counter++

	return o, nil
}

func (r *InMemoryOrderRepository) GetAll() []model.Order {

	responseOrders := make([]model.Order, len(r.orders))
	copy(responseOrders, r.orders)
	return responseOrders
}

func (r *InMemoryOrderRepository) GetOrderByID(id string) (*model.Order, error) {

	for o := range r.orders {
		if strings.EqualFold(r.orders[o].ID, id) {
			return &r.orders[o], nil
		}
	}
	return nil, fmt.Errorf("No order found for order id %s", id)
}

func (r *InMemoryOrderRepository) DeleteOrderByID(id string) (string, error) {

	for o := range r.orders {
		if strings.EqualFold(r.orders[o].ID, id) {

			r.orders = append(r.orders[:o], r.orders[o+1:]...)
			return fmt.Sprintf("Order with order id %s is deleted.", id), nil
		}
	}
	return "", fmt.Errorf("No order found for order id %s", id)
}

func (r *InMemoryOrderRepository) UpdateOrderByID(id string, o *model.Order) (*model.Order, error) {

	for index := range r.orders {
		if strings.EqualFold(r.orders[index].ID, id) {
			r.orders[index].CustomerName = o.CustomerName
			r.orders[index].Amount = o.Amount

			return &r.orders[index], nil

		}
	}

	return nil, fmt.Errorf("No order found with id %s", id)
}
