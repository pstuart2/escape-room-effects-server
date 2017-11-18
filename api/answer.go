package api

import (
	"github.com/labstack/echo"
	"net/http"
	"escape-room-effects-server/sound"
)

type AnswerRequest struct {
	Result string `json:"result"`
}

func (s Server) Answer(ctx echo.Context) error {
	r := new(AnswerRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}

	switch r.Result {
	case "correct":
		{
			go func() {
				sound.Play(sound.CorrectAnswer)
			}()
		}

	case "wrong":
		{
			playWrongAnswerSound()
		}

	}

	return nil
}

func playWrongAnswerSound() {
	go func() {
		sound.Play(sound.WrongAnswer)
	}()
}