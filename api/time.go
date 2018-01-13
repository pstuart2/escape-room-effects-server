package api

import (
	"github.com/labstack/echo"
	"net/http"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"math"
	"time"
	"gopkg.in/mgo.v2"
	"escape-room-effects-server/piClient"
)

type TimeRequest struct {
	Pos int `json:"pos"`
}

func (s *Server) Hours(ctx echo.Context) error {
	r := new(TimeRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	fmt.Printf("Hours: %d\n", r.Pos)

	db := s.getDb()
	defer db.Close()

	// Update db
	c := getGameCollection(db)
	c.UpdateId(s.GameID, bson.M{"$set": bson.M{"time.clock.hour": r.Pos}})

	return ctx.String(http.StatusOK, "")
}

func (s *Server) Minutes(ctx echo.Context) error {
	r := new(TimeRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	fmt.Printf("Minutes: %d\n", r.Pos)

	db := s.getDb()
	defer db.Close()

	// Update db
	c := getGameCollection(db)
	c.UpdateId(s.GameID, bson.M{"$set": bson.M{"time.clock.min": clockPosToMinute(r.Pos)}})

	return ctx.String(http.StatusOK, "")
}

func clockPosToMinute(pos int) int {
	return int(math.Ceil(float64(float64(pos) * 2.5)))
}

func isClockSetAheadByNoMoreThan5Min(s *Server, db *mgo.Session) bool {
	game := s.getGame(db)

	gt := game.Time.Current
	clockTime := time.Date(gt.Year(), gt.Month(), gt.Day(), game.Time.Clock.Hour, game.Time.Clock.Min, gt.Second(), gt.Nanosecond(), time.Local)

	diff := clockTime.Sub(gt)
	minutes := diff.Minutes()

	if minutes > 5 || minutes <= 0 {
		return false
	}

	return true
}

func (s *Server) StartTicker() {
	s.Ticker = time.NewTicker(time.Second * 1)
	fmt.Println("Ticker starting!")

	go func() {
		db := s.getDb()
		defer db.Close()

		lightsOn := -1  // Not yet set

		// TODO: Setup pi Zero for the lights server

		for range s.Ticker.C {
			// Check clock against time
			if isClockSetAheadByNoMoreThan5Min(s, db) {
				if lightsOn != 1 {
					fmt.Println("Time good, turning lights on")
					piClient.GameRoomLightsOnly()
					lightsOn = 1
				}
			} else if lightsOn != 0 {
				fmt.Println("Time fail, turning lights off")
				piClient.LightsOff()
				lightsOn = 0
			}
		}
	}()
}

func (s *Server) StopTicker() {
	if s.Ticker != nil {
		s.Ticker.Stop()
	}

	s.Ticker = nil
}
