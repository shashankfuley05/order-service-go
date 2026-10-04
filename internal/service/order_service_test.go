package service

import (
	"errors"
	"testing"

	"github.com/shashank/order-service/internal/dto"
	"github.com/shashank/order-service/internal/model"
)

type FakeOrderRepository struct {
	SaveFunc func(o *model.Order) (*model.Order, error)

	GetAllFunc func() []model.Order

	GetOrderByIDFunc func(id string) (*model.Order, error)

	DeleteOrderByIDFunc func(id string) (string, error)

	UpdateOrderByIDFunc func(id string, o *model.Order) (*model.Order, error)
}

func (f *FakeOrderRepository) Save(o *model.Order) (*model.Order, error) {
	return f.SaveFunc(o)
}

func (f *FakeOrderRepository) GetAll() []model.Order {
	return f.GetAllFunc()
}

func (f *FakeOrderRepository) GetOrderByID(id string) (*model.Order, error) {
	return f.GetOrderByIDFunc(id)
}

func (f *FakeOrderRepository) DeleteOrderByID(id string) (string, error) {
	return f.DeleteOrderByIDFunc(id)
}

func (f *FakeOrderRepository) UpdateOrderByID(id string, o *model.Order) (*model.Order, error) {
	return f.UpdateOrderByIDFunc(id, o)
}
func TestCreateOrder(t *testing.T) {

	fakeRepository := &FakeOrderRepository{
		SaveFunc: func(o *model.Order) (*model.Order, error) {
			o.ID = "TEST_ORD-1"
			return o, nil
		},
	}

	orderService := NewOrderService(fakeRepository)

	orderRequest := &dto.CreateOrderRequest{
		CustomerName: "The Ultimate Shashank Fuley",
		Amount:       38487800000009.00,
	}

	orderResponse, err := orderService.CreateOrder(orderRequest)

	if err != nil {
		t.Fatalf("BC error aagya %v", err)
	}

	if orderResponse == nil {
		t.Fatalf("BC response kaise nil aagya !!!")
	}

	if orderResponse.CustomerName != "The Ultimate Shashank Fuley" {
		t.Errorf("BC aur kaun ultimate aagya tere alawa..??")
	}

	if orderResponse.ID != "TEST_ORD-1" {
		t.Errorf("BC ye kya orders id aagayi %s", orderResponse.ID)
	}

}

func TestCreateOrder_Testing(t *testing.T) {

	tests := []struct {
		name         string
		amount       float64
		expectedErr  bool
		saveFunction func(o *model.Order) (*model.Order, error)
	}{
		{name: "Order Creation Successful",
			amount:      752056789468964.00,
			expectedErr: false,
			saveFunction: func(o *model.Order) (*model.Order, error) {
				return o, nil
			},
		},
		{name: "Invalid Amount",
			amount:      75205.00,
			expectedErr: true,
			saveFunction: func(o *model.Order) (*model.Order, error) {
				return nil, nil
			},
		},
		{name: "Repository Error",
			amount:      7520500000000.00,
			expectedErr: true,
			saveFunction: func(o *model.Order) (*model.Order, error) {
				return nil, errors.New("Repository ka error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeOrderRepository := &FakeOrderRepository{
				SaveFunc: tt.saveFunction,
			}

			orderService := NewOrderService(fakeOrderRepository)

			_, errr := orderService.CreateOrder(&dto.CreateOrderRequest{
				CustomerName: "The Ultimate legend Shashank Fuley",
				Amount:       tt.amount,
			})

			if tt.expectedErr == false {
				if errr != nil {
					t.Errorf("Bhai ye error kahaan se aagya")
				}
			} else {
				if errr == nil {
					t.Errorf("Bhai error aana chahiye tha")
				}
			}

		})
	}
}

func TestValidTransition(t *testing.T) {

	tests := []struct {
		name           string
		existingStatus string
		currentStatus  string
		expectedResult bool
	}{
		{
			name:           "CreatedToConfirmed",
			existingStatus: "CREATED",
			currentStatus:  "CONFIRMED",
			expectedResult: true,
		},
		{
			name:           "ConfirmedToShipped",
			existingStatus: "CONFIRMED",
			currentStatus:  "SHIPPED",
			expectedResult: true,
		},
		{
			name:           "ShippedToDelivered",
			existingStatus: "SHIPPED",
			currentStatus:  "DELIVERED",
			expectedResult: true,
		},
		{
			name:           "CreatedToShipped",
			existingStatus: "CREATED",
			currentStatus:  "SHIPPED",
			expectedResult: false,
		},
		{
			name:           "CreatedToDelivered",
			existingStatus: "CREATED",
			currentStatus:  "DELIVERED",
			expectedResult: false,
		},
		{
			name:           "DeliveredToConfirmed",
			existingStatus: "DELIVERED",
			currentStatus:  "CONFIRMED",
			expectedResult: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isValidTransition(test.existingStatus, test.currentStatus)
			if result != test.expectedResult {
				t.Errorf("Test failed for test %s, expected %v but got %v", test.name, test.expectedResult, result)
			}
		})
	}

}
