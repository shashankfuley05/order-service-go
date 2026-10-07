package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("%v", err)
	}

	log.Println("Server gracefully stopped")

}
