package api

import (
	"escape-room-effects-server/sound"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type AnswerRequest struct {
	Result string `json:"result"`
}

func (s *Server) Ask(game *GameState, db *mgo.Session) {
	current := getCurrentQuestion(game)
	if current == nil {
		current = getNextQuestion(game)
	}

	if current == nil {
		return
	}

	current.Asked = true

	c := getGameCollection(db)
	c.Update(bson.M{"_id": s.GameID, "questions.question": current.Query}, bson.M{
		"$set": bson.M{
			"say": current.Query,
			"questions.$": current,
		},
	})
}

func (s *Server) Answer(game *GameState, db *mgo.Session, answer string) {
	current := getCurrentQuestion(game)
	if current == nil {
		return
	}

	if strings.Contains(answer, current.Answer) {
		playCorrectSound()

		current.Answered = true

		c := getGameCollection(db)
		c.Update(bson.M{"_id": s.GameID, "questions.question": current.Query}, bson.M{
			"$set": bson.M{
				"say": current.Query,
				"questions.$": current,
			},
		})
	} else {
		playWrongAnswerSound()
	}
}

func getCurrentQuestion(game *GameState) (*Question) {
	for _, q := range game.Questions {
		if q.Asked && !q.Answered {
			return &q
		}
	}

	return nil
}

func getNextQuestion(game *GameState) (*Question) {
	for _, q := range game.Questions {
		if !q.Asked {
			return &q
		}
	}

	return nil
}

func playWrongAnswerSound() {
	go func() {
		sound.Play(sound.WrongAnswer)
	}()
}

func playCorrectSound() {
	go func() {
		sound.Play(sound.CorrectAnswer)
	}()
}
