package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wariusagi/assessment-tax/pkg/config"
	"github.com/wariusagi/assessment-tax/pkg/database"
	"github.com/wariusagi/assessment-tax/pkg/handlers"
	"github.com/wariusagi/assessment-tax/pkg/services"
)

func main() {
	config := config.NewConfig()
	db, err := database.InitDB(config.DatabaseUrl)
	if err != nil {
		log.Fatalf("Initialize database failed: %v", err)
	}
	defer database.CloseDB(db)

	e := setUpRoute(db)

	startServer(e, config.Port)
}

func setUpRoute(db *sql.DB) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	repo := database.NewRepositoryDB(db)
	taxService := services.NewTax(repo)
	taxHandler := handlers.NewTaxHandler(taxService)
	e.POST("/tax/calculations", taxHandler.CalculateTax)
	return e
}

func startServer(e *echo.Echo, port string) {
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	fmt.Println("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatalf("Error shutting down the server: %v", err)
	}
}
