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
			c.UpdateId(s.GameID, bson.M{"$set": bson.M{"camera": bson.M{"listeningState": Listening, "text": "Listening..."}}})
		}

	case ":getting-speech":
		{
			c.UpdateId(s.GameID, bson.M{"$set": bson.M{"camera": bson.M{"listeningState": Fetching, "text": "Decoding..."}}})
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
	fmt.Printf("Recorded: %s\n", r.Text)

	c := getGameCollection(db)
	game := s.getGame(db)

	if game.Eyes.State == Found && (strings.Contains(r.Text, "your name is") || strings.Contains(r.Text, "you are")) {
		return processShutdownCommand(s, ctx, r, db)
	} else if game.Eyes.State != Found && strings.Contains(r.Text, "do not be afraid") {
		c.UpdateId(s.GameID, bson.M{"$set": bson.M{"eyes": bson.M{"state": 0, "interact": Found}}})
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
