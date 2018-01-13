package api

import (
	"github.com/labstack/echo"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/labstack/gommon/log"
)

type FaceUpdateRequest struct {
	PreviousCount int    `json:"previousCount"`
	CurrentCount  int    `json:"currentCount"`
}

const (
	Hiding  = 0
	Waiting = 1
	Peeking = 2
	Found   = 3
	AfraidOfDark= 4
)

func (s *Server) Faces(ctx echo.Context) error {
	r := new(FaceUpdateRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	db := s.getDb()
	defer db.Close()

	game := s.getGame(db)

	if game.Eyes.Interact < Found {
		c := getGameCollection(db)

		if r.CurrentCount == 0 && game.Eyes.Interact == Hiding {
			log.Print("Setting Waiting")
			c.UpdateId(s.GameID, bson.M{"$set": bson.M{"eyes.interact": Waiting}})
		} else if game.Eyes.Interact != Hiding {
			log.Print("Setting Hiding")
			c.UpdateId(s.GameID, bson.M{"$set": bson.M{"eyes.interact": Hiding}})
		}
	}

	return ctx.String(http.StatusOK, "")
}
