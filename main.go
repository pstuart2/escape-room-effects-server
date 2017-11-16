package main

import (
	"github.com/labstack/echo"

	"os"
	"time"
	"os/signal"
	"context"
	"escape-room-effects-server/api"
)

const (
	PiServer = "http://192.168.86.101:8080"
)

func main() {
	startServer()
}

func startServer() {
	e := echo.New()

	e.POST("/state", api.GameState)
	e.POST("/answer", api.Answer)

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	api.StopRandomEffects()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
