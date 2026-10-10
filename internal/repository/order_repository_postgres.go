package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	apperrors "github.com/shashank/order-service/internal/errors"
	"github.com/shashank/order-service/internal/model"
)

type PostgresDB interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type PostgresOrderRepository struct {
	db PostgresDB
}

func NewPostgresOrderRepository(db PostgresDB) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		db: db,
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

	query := `SELECT id, customer_name, amount, status from orders`

	rows, err := p.db.Query(context.Background(), query)

	if err != nil {
		log.Printf("Error occured while fethcing orders %v", err)
		return []model.Order{}
	}

	defer rows.Close()

	reponse := make([]model.Order, 0)

	for rows.Next() {
		var tempModel model.Order

		err := rows.Scan(
			&tempModel.ID,
			&tempModel.CustomerName,
			&tempModel.Amount,
			&tempModel.Status,
		)

		if err != nil {
			return []model.Order{}
		}

		reponse = append(reponse, tempModel)
	}

	return reponse
}

func (p *PostgresOrderRepository) GetOrderByID(id string) (*model.Order, error) {

	query := `SELECT id, customer_name, amount, status from orders WHERE id = $1`

	row := p.db.QueryRow(context.Background(), query, id)

	var order model.Order

	err := row.Scan(
		&order.ID,
		&order.CustomerName,
		&order.Amount,
		&order.Status,
	)

	if err != nil {
		log.Printf("Error while fetching order by order id %s cause %v", id, err)
		return nil, apperrors.ErrOrderNotFound
	}

	return &order, nil
}

func (p *PostgresOrderRepository) DeleteOrderByID(id string) (string, error) {
	query := `DELETE FROM orders where id = $1`

	row, err := p.db.Exec(context.Background(), query, id)

	if err != nil {
		return "", err
	}

	if row.RowsAffected() == 0 {
		return "", apperrors.ErrOrderNotFound
	}

	return id, nil
}

func (p *PostgresOrderRepository) UpdateOrderByID(id string, o *model.Order) (*model.Order, error) {

	var udpatedOrder model.Order

	query := `UPDATE orders SET status = $1 WHERE id = $2 RETURNING id, customer_name, amount, status;`

	row := p.db.QueryRow(context.Background(), query, o.Status, id)

	err := row.Scan(
		&udpatedOrder.ID,
		&udpatedOrder.CustomerName,
		&udpatedOrder.Amount,
		&udpatedOrder.Status,
	)

	if err != nil {
		return nil, err
	}
	return &udpatedOrder, nil
}
