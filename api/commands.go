package api

import (
	"escape-room-effects-server/piClient"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CommandRequest struct {
	Command string `json:"command"`
}

// Command Processes commands from the game
func (s *Server) Command(ctx echo.Context) error {
	r := new(CommandRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}

	db := s.getDb()
	defer db.Close()

	fmt.Printf("Command [%s]\n", r.Command)

	c := getGameCollection(db)
	c.UpdateId(s.GameID, bson.M{"$push": bson.M{"commandsSent": bson.M{"command": r.Command}}})

	if s.isGamePaused(db) {
		return processPausedCommand(s, ctx, r, db)
	}

	return processCommand(s, ctx, r, db)
}

func processCommand(s *Server, ctx echo.Context, r *CommandRequest, db *mgo.Session) error {
	if strings.Contains(r.Command, "shutdown code") {
		return processShutdownCommand(s, ctx, r, db)
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
			s.pauseGame(db)
		}

	default:
		{
			return ctx.JSON(http.StatusBadRequest, "invalid command")
		}
	}

	return ctx.JSON(http.StatusOK, "OK")
}

func processPausedCommand(s *Server, ctx echo.Context, r *CommandRequest, db *mgo.Session) error {
	switch r.Command {
	case "resume game":
		{
			s.resumeGame(db)
		}

	default:
		{
			return ctx.JSON(http.StatusBadRequest, "invalid command while paused")
		}
	}

	return ctx.JSON(http.StatusOK, "OK")
}

func processShutdownCommand(s *Server, ctx echo.Context, r *CommandRequest, db *mgo.Session) error {
	shutdownCode := s.getShutdownCode(db)

	if r.Command != "shutdown code "+shutdownCode {
		playWrongAnswerSound()
		return ctx.JSON(http.StatusNotAcceptable, "Invalid shutdown code!")
	}

	s.finishGame(db)
	return ctx.JSON(http.StatusOK, "OK")
}
