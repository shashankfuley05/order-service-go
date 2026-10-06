package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v5"
	apperrors "github.com/shashank/order-service/internal/errors"
	"github.com/shashank/order-service/internal/model"
)

func TestUpdateOrderId(t *testing.T) {

	mockedPool, err := pgxmock.NewPool()

	if err != nil {
		t.Errorf("Can't create mock pool %v ", err)
	}

	defer mockedPool.Close()

	repository := NewPostgresOrderRepository(mockedPool)

	orderId := "ORD-1"
	status := "CONFIRMED"
	query := `UPDATE orders SET status = $1 WHERE id = $2 RETURNING id, customer_name, amount, status;`

	mockedPool.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(status, orderId).WillReturnRows(mockedPool.NewRows(
		[]string{"id", "customer_name", "amount", "status"},
	).AddRow("ORD-1", "The Great Shashank Fuley", 100000000.00, "CONFIRMED"))

	updatedOrder, err := repository.UpdateOrderByID("ORD-1", &model.Order{
		Status: status,
	})

	if err != nil {
		t.Fatalf("Got error %v", err)
	}

	if updatedOrder.ID != "ORD-1" {
		t.Errorf("Expected order id %s but got %s", orderId, updatedOrder.ID)
	}

	if updatedOrder.Status != status {
		t.Errorf("Expected status %s but got %s", status, updatedOrder.Status)
	}

	if err := mockedPool.ExpectationsWereMet(); err != nil {
		t.Fatalf("Mock expectations were not met %v", err)
	}
}

func TestGetOrderByIdForSuccess(t *testing.T) {

	mockedPool, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Can't create mocked pool %v", err)
	}
	defer mockedPool.Close()

	repository := NewPostgresOrderRepository(mockedPool)

	query := `SELECT id, customer_name, amount, status from orders WHERE id = $1`
	orderID := "ORD-1"
	customerName := "The Great Shashank Fuley"
	amount := 100000000.00
	status := "CONFIRMED"
	mockedPool.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(orderID).WillReturnRows(pgxmock.NewRows([]string{"id", "customer_name", "amount", "status"}).AddRow(orderID, customerName, amount, status))
	order, err := repository.GetOrderByID(orderID)

	if err != nil {
		t.Errorf("Error not expected but got %v", err)
	}

	if order.ID != orderID {
		t.Errorf("Expected order id was %s but got %s", orderID, order.ID)
	}

	if order.CustomerName != customerName {
		t.Errorf("Expected customer name is %s but got %s", customerName, order.CustomerName)
	}

	if order.Amount != amount {
		t.Errorf("Expected amount is %f but got %f", amount, order.Amount)
	}

	if order.Status != status {
		t.Errorf("Expected status is %s is but got %s", status, order.Status)
	}

	if err := mockedPool.ExpectationsWereMet(); err != nil {
		t.Fatalf("Mocked expectations were not met %v", err)
	}

}

func TestGetOrderByIdForException(t *testing.T) {

	mockedPool, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Can't create mocked pool %v", err)
	}
	defer mockedPool.Close()

	repository := NewPostgresOrderRepository(mockedPool)

	query := `SELECT id, customer_name, amount, status from orders WHERE id = $1`
	orderID := "ORD-1"

	mockedPool.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(orderID).WillReturnError(pgx.ErrNoRows)
	order, err := repository.GetOrderByID(orderID)

	if err == nil {
		t.Errorf("Error expected but didn't receive")
	}

	if !errors.Is(err, apperrors.ErrOrderNotFound) {

	}

	if order != nil {
		t.Errorf("Order was not expected but received")
	}

	if err := mockedPool.ExpectationsWereMet(); err != nil {
		t.Fatalf("Mocked expectations were not met %v", err)
	}

}

func TestDeleteOrderById(t *testing.T) {

	mockedPool, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Can't create mocked pool %v", err)
	}

	query := `DELETE FROM orders where id = $1`

	orderID := "ORD-1"

	mockedPool.ExpectExec(regexp.QuoteMeta(query)).WithArgs(orderID).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	repository := NewPostgresOrderRepository(mockedPool)

	message, err := repository.DeleteOrderByID(orderID)

	if err != nil {
		t.Errorf("Error was not expected but got %v", err)
	}

	if message != orderID {
		t.Errorf("Expected message was %s but got %s", orderID, message)
	}

	if err := mockedPool.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectation were not met %v", err)
	}
}

func TestDeleteOrderByIdForException(t *testing.T) {

	mockedPool, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Can't create mocked pool %v", err)
	}

	query := `DELETE FROM orders where id = $1`

	orderID := "ORD-1"

	mockedPool.ExpectExec(regexp.QuoteMeta(query)).WithArgs(orderID).WillReturnError(pgx.ErrTxClosed)

	repository := NewPostgresOrderRepository(mockedPool)

	message, err := repository.DeleteOrderByID(orderID)

	if err == nil {
		t.Errorf("Expected error but didn't recieve")
	}

	if message != "" {
		t.Errorf("Expected message was %s but got %s", "", message)
	}

	if err := mockedPool.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met %v", err)
	}

}

func TestDeleteOrderByIdForNowRowEffected(t *testing.T) {

	mockedPool, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Can't create mocked pool %v", err)
	}

	defer mockedPool.Close()

	query := `DELETE FROM orders where id = $1`

	orderID := "ORD-1"

	mockedPool.ExpectExec(regexp.QuoteMeta(query)).WithArgs(orderID).WillReturnResult(pgxmock.NewResult("DELETE", 0))

	repository := NewPostgresOrderRepository(mockedPool)

	message, err := repository.DeleteOrderByID(orderID)

	if err == nil {
		t.Errorf("Error was expected")
	}

	if !errors.Is(err, apperrors.ErrOrderNotFound) {
		t.Errorf("Expected error was of type %v but got %v", apperrors.ErrOrderNotFound, err)
	}

	if message != "" {
		t.Errorf("Expected message was %s but got %s", "", message)
	}

	if err := mockedPool.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met %v", err)
	}

}

func TestGetAllOrdersForSuccess(t *testing.T) {

	mockedPool, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Can't create mocked pool %v", err)
	}

	defer mockedPool.Close()

	query := `SELECT id, customer_name, amount, status from orders`

	mockedPool.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(
		pgxmock.NewRows([]string{"id", "customer_name", "amount", "status"}).
			AddRow("ORD-1", "The great Shashank Fuley", 324528732423.00, "CONFIRMED").
			AddRow("ORD-2", "Jim Moriarty", 9434399494334.00, "DELIVERED"),
	)

	repository := NewPostgresOrderRepository(mockedPool)

	orders := repository.GetAll()

	if len(orders) != 2 {
		t.Errorf("Expected number of orders were %v, but got %v", 2, len(orders))
	}

	if orders[0].ID != "ORD-1" {
		t.Errorf("Expected Order id was %s, but got %s", "ORD-1", orders[0].ID)
	}

	if orders[1].ID != "ORD-2" {
		t.Errorf("Expected Order id was %s, but got %s", "ORD-2", orders[1].ID)
	}

	if err := mockedPool.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met %v", err)
	}

}

func TestGetAllForEmptySliceIfDBException(t *testing.T) {

	mockedPool, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Can't create mocked pool %v", err)
	}

	defer mockedPool.Close()

	query := `SELECT id, customer_name, amount, status from orders`

	mockedPool.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(pgx.ErrNoRows)

	repository := NewPostgresOrderRepository(mockedPool)

	orders := repository.GetAll()

	if len(orders) > 0 {
		t.Errorf("Expectes length of slice was %d, but got %d", 0, len(orders))
	}

	if err := mockedPool.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met %v", err)
	}
}

func TestSaveOrderForSuccess(t *testing.T) {

	mockedPool, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Can't create mocked pool %v", err)
	}

	defer mockedPool.Close()

	query := ` INSERT INTO orders (customer_name, amount, status) 
			   VALUES ($1,$2,$3)
			   RETURNING id
	`
	order := &model.Order{
		CustomerName: "The Great Shashank Fuley",
		Amount:       9999999999.00,
	}
	mockedPool.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(order.CustomerName, order.Amount, "CREATED").WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("ORD-1"))

	repository := NewPostgresOrderRepository(mockedPool)

	order, err = repository.Save(order)

	if err != nil {
		t.Errorf("Error was not expected but got %v", err)
	}

	if order.ID != "ORD-1" {
		t.Errorf("Expectd order id was %s but got %s", "ORD-1", order.ID)
	}

	if order.Status != "CREATED" {
		t.Errorf("Expected status was %s but got %s", "CREATED", order.Status)
	}

	if err := mockedPool.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met %v", err)
	}
}

func TestSaveOrderForDBError(t *testing.T) {

	mockedPool, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Can't create mocked pool %v", err)
	}

	defer mockedPool.Close()

	query := ` INSERT INTO orders (customer_name, amount, status) 
			   VALUES ($1,$2,$3)
			   RETURNING id
	`
	order := &model.Order{
		CustomerName: "The Great Shashank Fuley",
		Amount:       9999999999.00,
	}
	mockedPool.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(order.CustomerName, order.Amount, "CREATED").WillReturnError(pgx.ErrNoRows)

	repository := NewPostgresOrderRepository(mockedPool)

	order, err = repository.Save(order)

	if err == nil {
		t.Errorf("Error was expected but got %v", err)
	}

	if order != nil {
		t.Errorf("Ordewr was expected to be nil")
	}

	if err := mockedPool.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met %v", err)
	}
}
