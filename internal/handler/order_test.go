package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shashank/order-service/internal/dto"
	apperrors "github.com/shashank/order-service/internal/errors"

	"github.com/shashank/order-service/internal/model"
)

type FakeOrderService struct {
	CreateOrderFunc     func(request *dto.CreateOrderRequest) (*model.Order, error)
	FetchOrdersFunc     func() []model.Order
	GetOrderByIDFunc    func(id string) (*model.Order, error)
	DeleteOrderByIDFunc func(id string) (string, error)
	UpdateOrderByIDFunc func(id string, request *dto.UpdateOrderRequest) (*model.Order, error)
}

func (s *FakeOrderService) CreateOrder(request *dto.CreateOrderRequest) (*model.Order, error) {
	return s.CreateOrderFunc(request)
}

func (s *FakeOrderService) FetchOrders() []model.Order {
	return s.FetchOrdersFunc()
}

func (s *FakeOrderService) GetOrderByID(id string) (*model.Order, error) {
	return s.GetOrderByIDFunc(id)
}

func (s *FakeOrderService) DeleteOrderByID(id string) (string, error) {
	return s.DeleteOrderByIDFunc(id)
}

func (s *FakeOrderService) UpdateOrderByID(id string, request *dto.UpdateOrderRequest) (*model.Order, error) {
	return s.UpdateOrderByIDFunc(id, request)
}
func TestCreateOrderHandler(t *testing.T) {

	tests := []struct {
		name           string
		expectedStatus int
		orderRequest   *dto.CreateOrderRequest
	}{
		{
			name:           "Invalid Amount",
			expectedStatus: http.StatusBadRequest,
			orderRequest: &dto.CreateOrderRequest{
				CustomerName: "tom hardy",
				Amount:       4343254.00,
			},
		},
		{
			name:           "Order Created Successfully",
			expectedStatus: http.StatusCreated,
			orderRequest: &dto.CreateOrderRequest{
				CustomerName: "tom hardy",
				Amount:       432654636634.00,
			},
		},
	}

	fakeOrderService := &FakeOrderService{}
	fakeOrderService.CreateOrderFunc = func(request *dto.CreateOrderRequest) (*model.Order, error) {
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
		order.Status = "TEST_STATUS_CREATED"
		order.ID = "TEST_ID"

		return order, nil

	}

	orderHadler := NewOrderHandler(fakeOrderService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/orders", orderHadler.CreateOrder)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			requestByte, _ := json.Marshal(test.orderRequest)
			req := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(requestByte))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			if test.expectedStatus != w.Code {
				t.Errorf("expected status %d, got %d", test.expectedStatus, w.Code)
			}

			if test.expectedStatus == http.StatusCreated {
				var response dto.CreateOrderResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}

				if response.ID != "TEST_ID" {
					t.Errorf("Order status mismatch expected %s and got %s", "TEST_ID", response.ID)
				}

			}

			if test.expectedStatus == http.StatusBadRequest {
				var response map[string]string

				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}

				if response["error"] != "order validation failed : Amount is invalid" {
					t.Errorf("expected error message %s, got %s",
						"order validation failed : Amount is invalid",
						response["error"])
				}
			}
		})
	}

}

func TestGetOrderByID(t *testing.T) {

	tests := []struct {
		name           string
		expectedStatus int
		orderId        string
	}{
		{
			name:           "Order not found with id",
			expectedStatus: http.StatusNotFound,
			orderId:        "ORDER_ID_NOTFOUND",
		},
		{
			name:           "Order found",
			expectedStatus: http.StatusOK,
			orderId:        "ORDER_ID_FOUND",
		},
	}

	orderService := &FakeOrderService{}

	orderService.GetOrderByIDFunc = func(id string) (*model.Order, error) {
		if id == "ORDER_ID_NOTFOUND" {
			return nil, fmt.Errorf("Order not found with id %s : %w", id, apperrors.ErrOrderNotFound)
		}
		return &model.Order{
			ID: "ORDER_ID_FOUND",
		}, nil
	}

	orderHandler := NewOrderHandler(orderService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/orders/:id", orderHandler.GetOrderByID)

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/orders/"+test.orderId, nil)

			router.ServeHTTP(w, request)

			if w.Code != test.expectedStatus {
				t.Errorf("Test faield expected status %d got %d", test.expectedStatus, w.Code)
			}

			if w.Code == http.StatusNotFound {

				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Could't parse the response %v", err)
				}

				if response["error"] != "Order not found with id ORDER_ID_NOTFOUND : Order not found" {
					t.Errorf("Expected error message got %s", response["error"])
				}
			}

			if w.Code == http.StatusOK {
				var response dto.CreateOrderResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Could't parse the response %v", err)
				}

				if response.ID != "ORDER_ID_FOUND" {
					t.Errorf("Expected error message got %s", response.ID)
				}
			}
		})

	}
}
