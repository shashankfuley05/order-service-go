package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shashank/order-service/internal/database"
	"github.com/shashank/order-service/internal/handler"
	"github.com/shashank/order-service/internal/middleware"
	"github.com/shashank/order-service/internal/repository"
	"github.com/shashank/order-service/internal/service"
)

func main() {

	pool, err := database.NewPostgresConnection()

	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	fmt.Println("PostgreSQL connected successfully")

	router := gin.Default()

	tokenHandler := handler.NewTokenHandler(&service.AuthServiceImpl{})

	router.Use(middleware.Logger())

	router.POST("/token", tokenHandler.GenerateTokenForUser)

	orderRepository := repository.NewPostgresOrderRepository(pool)

	orderService := service.NewOrderService(orderRepository)

	orderHandler := handler.NewOrderHandler(orderService)

	routerGroup := router.Group("")

	routerGroup.Use(middleware.AuthMiddleware())

	routerGroup.POST("/orders", middleware.RequiredRole([]string{"ADMIN"}), orderHandler.CreateOrder)

	routerGroup.GET("/orders", middleware.RequiredRole([]string{"USER", "ADMIN"}), orderHandler.FetchOrders)

	routerGroup.GET("/orders/:id", middleware.RequiredRole([]string{"USER", "ADMIN"}), orderHandler.GetOrderByID)

	routerGroup.DELETE("/orders/:id", middleware.RequiredRole([]string{"ADMIN"}), orderHandler.DeleteOrderByID)

	routerGroup.PUT("/orders/:id", middleware.RequiredRole([]string{"ADMIN"}), orderHandler.UpdateOrderByID)

	router.Run(":8080")

}
