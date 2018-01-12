package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	None      = 0
	Listening = 1
	Fetching  = 2
	Success   = 3
	Failed    = 4
)

type CommandRequest struct {
	Command string `json:"command"`
	Text    string `json:"text"`
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

	if strings.HasPrefix(r.Command, ":") {
		return processAppCommands(s, ctx, r, db)
	}

	return processSpokenCommand(s, ctx, r, db)
}

func processAppCommands(s *Server, ctx echo.Context, r *CommandRequest, db *mgo.Session) error {
	c := getGameCollection(db)

	switch r.Command {
	case ":listening":
		{
			c.UpdateId(s.GameID, bson.M{"$set": bson.M{"camera.listeningState": Listening}})
		}

	case ":getting-speech":
		{
			c.UpdateId(s.GameID, bson.M{"$set": bson.M{"camera.listeningState": Fetching}})
		}

	case ":speech":
		{
			c.UpdateId(s.GameID, bson.M{
				"$set": bson.M{
					"camera": bson.M{"listeningState": Success},
					"say":    r.Text,
				},
				"$push": bson.M{"recordings": r.Text}},
			)

			return processSpokenCommand(s, ctx, r, db)
		}

	case ":stopped":
		{
			// no-audio, could-not-translate, api-error, api-timeout
			c.UpdateId(s.GameID, bson.M{"$set": bson.M{"camera.listeningState": None, "say": getSayText(r.Text)}})
		}
	}

	return ctx.JSON(http.StatusOK, "OK")
}

func getSayText(text string) string {
	switch text {
	case "no-audio":
		return ""
	case "could-not-translate":
		return "I did not understand. Please speak loud and clear."
	}

	return "Something went wrong, please try again."
}

func processSpokenCommand(s *Server, ctx echo.Context, r *CommandRequest, db *mgo.Session) error {
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
