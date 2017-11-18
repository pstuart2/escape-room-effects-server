package api

import (
	"github.com/labstack/echo"
	"net/http"
	"escape-room-effects-server/piClient"
	"fmt"
	"gopkg.in/mgo.v2"
	"strings"
)

type CommandRequest struct {
	Command string `json:"command"`
}

func (s Server) Command(ctx echo.Context) error {
	r := new(CommandRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}

	db := s.getDb()
	defer db.Close()

	fmt.Printf("Command [%s]\n", r.Command)

	if isGamePaused(db) {
		return processPausedCommand(ctx, r, db)
	}

	return processCommand(ctx, r, db)
}

func processCommand(ctx echo.Context, r *CommandRequest, db *mgo.Session) error {
	if strings.Contains(r.Command, "shutdown code") {
		return processShutdownCommand(ctx, r, db)
	}

	switch r.Command {
	case "lights on":
		{
			piClient.GameRoomLightsOnly()
		}

	case "lights off":
		{
			piClient.LightsOff()
		}

	case "secret light":
		{
			piClient.SecretLight()
		}

	case "pause game":
		{
			pauseGame(db)
		}

	default:
		{
			return ctx.JSON(http.StatusBadRequest, "invalid command")
		}
	}

	return ctx.JSON(http.StatusOK, "OK")
}

func processPausedCommand(ctx echo.Context, r *CommandRequest, db *mgo.Session) error {
	switch r.Command {
	case "resume game":
		{
			resumeGame(db)
		}

	default:
		{
			return ctx.JSON(http.StatusBadRequest, "invalid command while paused")
		}
	}

	return ctx.JSON(http.StatusOK, "OK")
}

func processShutdownCommand(ctx echo.Context, r *CommandRequest, db *mgo.Session) error {
	shutdownCode := getShutdownCode(db)

	if r.Command != "shutdown code " + shutdownCode {
		playWrongAnswerSound()
		return ctx.JSON(http.StatusNotAcceptable, "Invalid shutdown code!")
	}

	finishGame(db)
	return ctx.JSON(http.StatusOK, "OK")
}