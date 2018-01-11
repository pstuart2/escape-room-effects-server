package api

import (
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
		return ctx.JSON(http.StatusOK, "OK")
	}

	return processCommand(s, ctx, r, db)
}

func processCommand(s *Server, ctx echo.Context, r *CommandRequest, db *mgo.Session) error {
	if strings.Contains(r.Command, "your name is") || strings.Contains(r.Command, "you are") {
		// TODO: These are invalid commands until they have me come out
		return processShutdownCommand(s, ctx, r, db)
	} else if strings.Contains(r.Command, "do not be afraid") {
		// TODO: Handle
	} else {
		fmt.Println("Invalid command")
	}

	return ctx.JSON(http.StatusOK, "OK")
}

func processShutdownCommand(s *Server, ctx echo.Context, r *CommandRequest, db *mgo.Session) error {
	shutdownCode := s.getShutdownCode(db)

	if strings.Contains(r.Command, shutdownCode) {
		s.finishGame(db)
	} else {
		playWrongAnswerSound()
	}

	return ctx.JSON(http.StatusOK, "OK")
}
