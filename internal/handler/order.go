package handler

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shashank/order-service/internal/dto"
	"github.com/shashank/order-service/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: s}
}

func (o *OrderHandler) CreateOrder(ctx *gin.Context) {
	var request dto.CreateOrderRequest
	error := ctx.ShouldBindJSON(&request)

	if error != nil {
		mapValidationError(ctx, error, request)
		return
	}

	order, error := o.orderService.CreateOrder(&request)

	if error != nil {
		writeError(ctx, http.StatusBadRequest, error.Error())
		return
	}

	response := dto.CreateOrderResponse{
		CustomerName: order.CustomerName,
		Amount:       order.Amount,
		ID:           order.ID,
		Status:       order.Status,
	}

	ctx.JSON(http.StatusCreated, response)

}

func (o *OrderHandler) FetchOrders(ctx *gin.Context) {

	responses := o.orderService.FetchOrders()
	finalResponse := make([]dto.CreateOrderResponse, 0, len(responses))

	for _, order := range responses {
		finalResponse = append(finalResponse, dto.CreateOrderResponse{
			CustomerName: order.CustomerName,
			Amount:       order.Amount,
			ID:           order.ID,
			Status:       order.Status,
		})
	}
	ctx.JSON(http.StatusOK, finalResponse)
}

func (o *OrderHandler) GetOrderByID(ctx *gin.Context) {
	response, error := o.orderService.GetOrderByID(ctx.Param("id"))

	if error != nil {
		writeError(ctx, http.StatusNotFound, error.Error())
		return
	}

	ctx.JSON(http.StatusOK, dto.CreateOrderResponse{
		CustomerName: response.CustomerName,
		Amount:       response.Amount,
		ID:           response.ID,
		Status:       response.Status,
	})

}

func (o *OrderHandler) DeleteOrderByID(ctx *gin.Context) {
	message, error := o.orderService.DeleteOrderByID(ctx.Param("id"))

	if error != nil {
		writeError(ctx, http.StatusNotFound, error.Error())
		return
	}

	ctx.JSON(http.StatusOK, message)
}

func (o *OrderHandler) UpdateOrderByID(ctx *gin.Context) {
	var request dto.UpdateOrderRequest
	err := ctx.ShouldBindJSON(&request)

	if err != nil {
		mapValidationError(ctx, err, request)
		return
	}

	order, err := o.orderService.UpdateOrderByID(ctx.Param("id"), &request)

	if err != nil {

		writeError(ctx, http.StatusNotFound, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, dto.CreateOrderResponse{
		CustomerName: order.CustomerName,
		Amount:       order.Amount,
		ID:           order.ID,
		Status:       order.Status,
	})
}

func writeError(ctx *gin.Context, status int, err string) {

	ctx.JSON(status, dto.ErrorResponse{
		Error: err,
	})

}

func mapValidationError(ctx *gin.Context, err error, request any) {

	t := reflect.TypeOf(request)

	validationErrors, ok := err.(validator.ValidationErrors)

	if ok {

		for _, validationError := range validationErrors {

			f, _ := t.FieldByName(validationError.Field())
			jsonTag := f.Tag.Get("json")
			switch validationError.Tag() {
			case "gt":
				writeError(ctx, http.StatusBadRequest, fmt.Sprintf("%s must be greater than 0", jsonTag))
				return

			case "required":
				writeError(ctx, http.StatusBadRequest, fmt.Sprintf("%s is required", jsonTag))
				return
			}

		}
	}

	writeError(ctx, http.StatusBadRequest, err.Error())

}
