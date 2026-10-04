package repository

import (
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock/v5"
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
