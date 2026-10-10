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
	"github.com/joho/godotenv"
	"github.com/shashank/order-service/internal/database"
	"github.com/shashank/order-service/internal/handler"
	"github.com/shashank/order-service/internal/middleware"
	"github.com/shashank/order-service/internal/repository"
	"github.com/shashank/order-service/internal/service"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Couldn't load environment %v", err)
	}

	postgresConfig := &database.PostgresConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
	}

	if err := validatePostgresConfig(postgresConfig); err != nil {
		log.Fatalf("Config not found %v", err)
	}

	pool, err := database.NewPostgresConnection(postgresConfig)

	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	log.Println("PostgreSQL connected successfully")

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
		Addr:    os.Getenv("SERVER_PORT"),
		Handler: router,
	}
	if err := validateServerConfig(server); err != nil {
		log.Fatalf("%v", err)
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
		log.Printf("Couldn't shutdown gracefully %v", err)
	}

}

func validateServerConfig(server *http.Server) error {
	if server.Addr == "" {
		return fmt.Errorf("SERVER_PORT is requireds")
	}

	return nil
}
func validatePostgresConfig(config *database.PostgresConfig) error {
	if config.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if config.Port == "" {
		return fmt.Errorf("DB_PORT is required")
	}
	if config.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if config.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if config.Database == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	return nil
}
