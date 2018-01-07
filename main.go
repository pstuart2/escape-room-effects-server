package main

import (
	"github.com/labstack/echo"

	"os"
	"time"
	"os/signal"
	"context"
	"escape-room-effects-server/api"
	"log"

	"gopkg.in/mgo.v2"
)

func main() {
	startServer()
}

func startServer() {
	e := echo.New()

	session := setupDatabase()
	defer session.Close()

	server := api.Server{Db: session}

	e.POST("/faces", server.Faces)
	e.POST("/state", server.GameState)
	e.POST("/answer", server.Answer)
	e.POST("/command", server.Command)
	e.POST("/hours", server.Hours)
	e.POST("/minutes", server.Minutes)

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

func setupDatabase() *mgo.Session {
	appDbMasterSession, err := mgo.Dial("mongodb://localhost:3001/meteor")
	if err != nil {
		log.Fatal("Failed to dial appDbMasterSession: " + err.Error())
	}

	appDbMasterSession.SetMode(mgo.Monotonic, true)

	return appDbMasterSession
}