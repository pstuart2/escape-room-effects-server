package api

import (
	"github.com/labstack/echo"
	"net/http"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"math"
	"time"
	"escape-room-effects-server/piClient"
	"escape-room-effects-server/sound"
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

func isClockSetAheadByNoMoreThan5Min(game *GameState) (bool, float64) {
	gt := game.Time.Current
	clockTime := time.Date(gt.Year(), gt.Month(), gt.Day(), game.Time.Clock.Hour, game.Time.Clock.Min, 0, gt.Nanosecond(), time.Local)

	diff := clockTime.Sub(gt)
	minutes := diff.Minutes()
	seconds := diff.Seconds()

	if minutes > 5 || minutes <= 0 {
		return false, seconds
	}

	return true, seconds
}

func (s *Server) StartTicker() {
	s.Ticker = time.NewTicker(time.Second * 1)
	fmt.Println("Ticker starting!")

	go func() {
		db := s.getDb()
		defer db.Close()

		lightsOn := -1 // Not yet set

		for range s.Ticker.C {
			game := s.getGame(db)

			// Check clock against time
			isAhead, seconds := isClockSetAheadByNoMoreThan5Min(game)
			if isAhead {
				if lightsOn != 1 {
					fmt.Println("Time good, turning lights on")
					piClient.GameRoomLightsOnly()
					lightsOn = 1
				} else {
					fmt.Printf("Seconds: %f\n", seconds)
					if seconds == 9 {
						sound.Play(sound.AliceIntro)
					}
				}
			} else if lightsOn != 0 {
				fmt.Println("Time fail, turning lights off")
				piClient.LightsOff()
				lightsOn = 0
			}

			if lightsOn == 0 && game.Eyes.Interact == Found {
				fmt.Println("Afraid of the dark!")
				c := getGameCollection(db)
				c.UpdateId(s.GameID, bson.M{"$set": bson.M{"eyes.interact": AfraidOfDark, "say": "I'm afraid of the dark."}})
			} else if lightsOn == 1 && game.Eyes.Interact == AfraidOfDark {
				c := getGameCollection(db)
				c.UpdateId(s.GameID, bson.M{"$set": bson.M{"eyes": bson.M{"interact": Found, "state": 0}, "say": "I'm glad it's not dark."}})
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
