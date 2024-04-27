package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wariusagi/assessment-tax/pkg/config"
	"github.com/wariusagi/assessment-tax/pkg/database"
)

func main() {
	config := config.NewConfig()

	e := echo.New()

	if err := database.InitDB(config.DatabaseUrl); err != nil {
		e.Logger.Fatalf("Initialize database failed: %v", err)
	}
	defer database.CloseDB()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	startServer(e, config.Port)
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
