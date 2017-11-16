package api

import (
	"github.com/labstack/echo"
	"net/http"
	"escape-room-effects-server/piClient"
	"fmt"
	"gopkg.in/mgo.v2"
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

	validCommand := false
	if isGamePaused(db) {
		validCommand = processPausedCommand(r, db)
	} else {
		validCommand = processCommand(r, db)
	}

	if !validCommand {
		return ctx.JSON(http.StatusBadRequest, "invalid command")
	}

	return ctx.JSON(http.StatusOK, "OK")
}

func processCommand(r *CommandRequest, db *mgo.Session) bool {
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
			return false
		}
	}

	return true
}

func processPausedCommand(r *CommandRequest, db *mgo.Session) bool {
	switch r.Command {
	case "resume game":
		{
			resumeGame(db)
		}

	default:
		{
			return false
		}
	}

	return true
}
