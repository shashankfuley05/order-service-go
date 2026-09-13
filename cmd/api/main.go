package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shashank/order-service/internal/handler"
	"github.com/shashank/order-service/internal/middleware"
	"github.com/shashank/order-service/internal/repository"
	"github.com/shashank/order-service/internal/service"
)

func main() {

	router := gin.Default()

	router.Use(middleware.Logger())

	router.Use(middleware.AuthMiddleware())

	orderRepository := repository.NewInMemoryOrderRepository()

	orderService := service.NewOrderService(orderRepository)

	orderHandler := handler.NewOrderHandler(orderService)

	router.POST("/orders", orderHandler.CreateOrder)

	router.GET("/orders", orderHandler.FetchOrders)

	router.GET("/orders/:id", orderHandler.GetOrderByID)

	router.DELETE("/orders/:id", orderHandler.DeleteOrderByID)

	router.PUT("/orders/:id", orderHandler.UpdateOrderByID)

	router.Run(":8080")
}
