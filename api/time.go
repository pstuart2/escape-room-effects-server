package api

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/labstack/gommon/log"
)

type TimeRequest struct {
	Pos int `json:"pos"`
}

func (s Server) Hours(ctx echo.Context) error {
	r := new(TimeRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	log.Printf("Hours: %d", r.Pos)

	//db := s.getDb()
	//defer db.Close()
	//
	//game := getGame(db)
	//
	//if game.Eyes.Interact != Found {
	//	c := getGameCollection(db)
	//
	//	if r.CurrentCount == 0 && game.Eyes.Interact == Hiding {
	//		log.Print("Setting Waiting")
	//		c.UpdateId(runningGameID, bson.M{"$set": bson.M{"eyes.interact": Waiting}})
	//	} else if game.Eyes.Interact != Hiding {
	//		log.Print("Setting Hiding")
	//		c.UpdateId(runningGameID, bson.M{"$set": bson.M{"eyes.interact": Hiding}})
	//	}
	//}

	return ctx.String(http.StatusOK, "")
}

func (s Server) Minutes(ctx echo.Context) error {
	r := new(TimeRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	log.Printf("Minutes: %d", r.Pos)

	return ctx.String(http.StatusOK, "")
}