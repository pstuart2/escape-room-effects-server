package api

import (
	"github.com/labstack/echo"
	"net/http"
	"escape-room-effects-server/sound"
	"escape-room-effects-server/piClient"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
var runningGameID string

type GameState struct {
	ID    string `bson:"_id"`
	State int    `bson:"state"`
}

func (s Server) GameState(ctx echo.Context) error {
	r := new(GameStateRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}

	runningGameID = getID(r)

	db := s.getDb()
	defer db.Close()

	process(r, db)

	return nil
}

func process(r *GameStateRequest, db *mgo.Session) {
	switch r.State {
	case "starting":
		{
			piClient.LightsOff()
			go func() {
				//sound.Play(sound.LightToggle)
				sound.Play(sound.DoorSlam)
				sound.Play(sound.ChainDoorShut)
			}()
		}

	case "start":
		{
			piClient.WallLightsOnly()
			startRandomEffects()
			go func() {
				sound.Play(sound.MusicLoop)
				sound.Play(sound.MusicLoop)
				sound.Play(sound.MusicLoop)
				sound.Play(sound.UndergroundEffect)
			}()
		}

	case "pause":
		{
			pauseGame(db)
		}

	case "resume":
		{
			resumeGame(db)
		}

	case "finish":
		{
			StopRandomEffects()
			piClient.LightsOn()
			go func() {
				sound.Play(sound.Clapping)
			}()
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

func getID(r *GameStateRequest) string {
	if len(r.ID) > 0 {
		return r.ID
	}

	return runningGameID
}

func resumeGame(db *mgo.Session) {
	c := getGameCollection(db)
	c.UpdateId(runningGameID, bson.M{"$set": bson.M{"state": Running}})

	piClient.WallLightsOnly()
	startRandomEffects()
	go func() {
		sound.Play(sound.Unpause)
	}()
}

func pauseGame(db *mgo.Session) {
	StopRandomEffects()

	c := getGameCollection(db)
	c.UpdateId(runningGameID, bson.M{"$set": bson.M{"state": Paused}, "$inc": bson.M{"timesPaused": 1}})

	piClient.LightsOff()
	go func() {
		sound.Play(sound.LightToggle)
	}()
}

func isGamePaused(db *mgo.Session) bool {
	c := getGameCollection(db)

	game := GameState{}

	if err := c.FindId(runningGameID).
		Select(bson.M{
		"_id":   1,
		"state": 1,
	}).One(&game); err != nil {
		return true
	}

	return game.State == Paused
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
