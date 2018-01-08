package api

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type Server struct {
	Db     *mgo.Session
	GameID string
	Ticker *time.Ticker
}

type Eyes struct {
	State    int `bson:"state"`
	Interact int `bson:"interact"`
}

type GameClock struct {
	Hour int `bson:"hour"`
	Min  int `bson:"min"`
}

type GameTime struct {
	Current time.Time `bson:"current"`
	Clock   GameClock `bson:"clock"`
}

type GameState struct {
	ID           string   `bson:"_id"`
	State        int      `bson:"state"`
	ShutdownCode string   `bson:"shutdownCode"`
	Eyes         Eyes     `bson:"eyes"`
	Time         GameTime `bson:"time"`
}

func (s *Server) getDb() *mgo.Session {
	return s.Db.Copy()
}

func getGameCollection(db *mgo.Session) *mgo.Collection {
	return db.DB("").C("game")
}

func (s *Server) isGamePaused(db *mgo.Session) bool {
	game := s.getGame(db)
	return game.State == Paused
}

func (s *Server) getShutdownCode(db *mgo.Session) string {
	game := s.getGame(db)
	return game.ShutdownCode
}

func (s *Server) getGame(db *mgo.Session) *GameState {
	c := getGameCollection(db)

	game := GameState{}

	if err := c.FindId(s.GameID).
		Select(bson.M{
		"_id":          1,
		"state":        1,
		"shutdownCode": 1,
		"eyes":         1,
		"time":         1,
	}).One(&game); err != nil {
		return nil
	}

	return &game
}
