package api

import (
	"github.com/labstack/echo"
	"net/http"
	"escape-room-effects-server/sound"
	"escape-room-effects-server/piClient"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

const (
	Pending  = 0
	Starting = 1
	Running  = 2
	Paused   = 3
	Finished = 4
)

type GameStateRequest struct {
	ID    string `json:"id"`
	State string `json:"state"`
}

var randomEffectChannel chan bool

func (s *Server) GameState(ctx echo.Context) error {
	r := new(GameStateRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}

	db := s.getDb()
	defer db.Close()

	process(s, r, db)

	return nil
}

func process(s *Server, r *GameStateRequest, db *mgo.Session) {
	fmt.Printf("Set State: %s", r.State)

	switch r.State {
	case "setGameId":
		{
			s.GameID = r.ID
			fmt.Printf("GameId set to: %s\n", s.GameID)
			game := s.getGame(db)

			if game.State == Running {
				s.StartTicker()
				startRandomEffects()
			}
		}
	case "starting":
		{
			s.GameID = r.ID

			piClient.WallLightsOnly()
			go func() {
				sound.Play(sound.DoorSlam)
				sound.Play(sound.ChainDoorShut)
				sound.Play(sound.Wonderland)
			}()
		}

	case "start":
		{
			s.StartTicker()
			piClient.LightsOff()
			startRandomEffects()
		}

	case "pause":
		{
			s.pauseGame(db)
		}

	case "resume":
		{
			s.resumeGame(db)
		}

	case "finish":
		{
			s.finishGame(db)
		}

	case "lightsOn":
		{
			piClient.LightsOn()
		}

	case "lightsOff":
		{
			piClient.LightsOff()
		}

	case "wallLightsOnly":
		{
			piClient.WallLightsOnly()
		}

	case "gameRoomLightsOnly":
		{
			piClient.GameRoomLightsOnly()
		}
	}
}

func (s *Server) resumeGame(db *mgo.Session) {
	s.StartTicker()
	c := getGameCollection(db)
	c.UpdateId(s.GameID, bson.M{"$set": bson.M{"state": Running}})

	piClient.WallLightsOnly()
	startRandomEffects()
	go func() {
		sound.Play(sound.Unpause)
	}()
}

func (s *Server) pauseGame(db *mgo.Session) {
	s.StopTicker()
	StopRandomEffects()

	c := getGameCollection(db)
	c.UpdateId(s.GameID, bson.M{"$set": bson.M{"state": Paused}, "$inc": bson.M{"timesPaused": 1}})

	piClient.LightsOff()
	go func() {
		sound.Play(sound.LightToggle)
	}()
}

func (s *Server) finishGame(db *mgo.Session) {
	s.StopTicker()
	c := getGameCollection(db)
	c.UpdateId(s.GameID, bson.M{"$set": bson.M{"state": Finished}})

	StopRandomEffects()
	piClient.LightsOn()
	go func() {
		sound.Play(sound.Clapping)
	}()
}

func startRandomEffects() {
	randomEffectChannel = sound.StartRandomEffects()
}

func StopRandomEffects() {
	if randomEffectChannel != nil {
		randomEffectChannel <- true
		randomEffectChannel = nil
	}
}
