package api

import (
	"github.com/labstack/echo"
	"net/http"
	"escape-room-effects-server/sound"
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
			go func() {
				//sound.Play(sound.LightToggle)
				sound.Play(sound.DoorSlam)
				sound.Play(sound.ChainDoorShut)
			}()
		}

	case "start":
		{
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
			go func() {
				sound.Play(sound.LightToggle)
			}()
		}

	case "resume":
		{
			startRandomEffects()
			go func() {
				sound.Play(sound.Unpause)
			}()
		}

	case "finish":
		{
			StopRandomEffects()
			go func() {
				sound.Play(sound.Clapping)
			}()
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