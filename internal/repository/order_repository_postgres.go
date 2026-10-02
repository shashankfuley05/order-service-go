package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shashank/order-service/internal/model"
)

type PostgresOrderRepository struct {
	db *pgxpool.Pool
}

func NewPostgresOrderRepository(pool *pgxpool.Pool) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		db: pool,
	}
}

func (p *PostgresOrderRepository) Save(o *model.Order) (*model.Order, error) {

	query := ` INSERT INTO orders (customer_name, amount, status) 
			   VALUES ($1,$2,$3)
			   RETURNING id
	`

	o.Status = "CREATED"
	err := p.db.QueryRow(context.Background(), query, o.CustomerName, o.Amount, o.Status).Scan(&o.ID)

	if err != nil {
		return nil, err
	}

	return o, nil
}

func (p *PostgresOrderRepository) GetAll() []model.Order {
	return nil
}

func (p *PostgresOrderRepository) GetOrderByID(id string) (*model.Order, error) {
	return nil, nil
}

func (p *PostgresOrderRepository) DeleteOrderByID(id string) (string, error) {
	return "", nil
}

func (p *PostgresOrderRepository) UpdateOrderByID(id string, o *model.Order) (*model.Order, error) {
	return nil, nil
}
