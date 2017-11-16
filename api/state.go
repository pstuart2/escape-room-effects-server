package api

import (
	"github.com/labstack/echo"
	"net/http"
	"escape-room-effects-server/sound"
	"escape-room-effects-server/piClient"
)

type GameStateRequest struct {
	State string `json:"state"`
}

var randomEffectChannel chan bool

func GameState(ctx echo.Context) error {
	r := new(GameStateRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}

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
			StopRandomEffects()
			piClient.LightsOff()
			go func() {
				sound.Play(sound.LightToggle)
			}()
		}

	case "resume":
		{
			piClient.WallLightsOnly()
			startRandomEffects()
			go func() {
				sound.Play(sound.Unpause)
			}()
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

	return nil
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