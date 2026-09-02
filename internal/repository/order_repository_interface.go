package repository

import "github.com/shashank/order-service/internal/model"

type OrderRepository interface {
	Save(o *model.Order) (*model.Order, error)
	GetAll() []model.Order
	GetOrderByID(id string) (*model.Order, error)
	DeleteOrderByID(id string) (string, error)
	UpdateOrderByID(id string, o *model.Order) (*model.Order, error)
}
