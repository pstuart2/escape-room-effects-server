package api

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type Server struct {
	Db *mgo.Session
}

type Eyes struct {
	State    int `bson:"state"`
	Interact int `bson:"interact"`
}

type GameState struct {
	ID           string `bson:"_id"`
	State        int    `bson:"state"`
	ShutdownCode string `bson:"shutdownCode"`
	Eyes         Eyes   `bson:"eyes"`
}

func (s Server) getDb() *mgo.Session {
	return s.Db.Copy()
}

func getGameCollection(db *mgo.Session) *mgo.Collection {
	return db.DB("").C("game")
}

func isGamePaused(db *mgo.Session) bool {
	game := getGame(db)
	return game.State == Paused
}

func getShutdownCode(db *mgo.Session) string {
	game := getGame(db)
	return game.ShutdownCode
}

func getGame(db *mgo.Session) *GameState {
	c := getGameCollection(db)

	game := GameState{}

	if err := c.FindId(runningGameID).
		Select(bson.M{
		"_id":          1,
		"state":        1,
		"shutdownCode": 1,
		"eyes":         1,
	}).One(&game); err != nil {
		return nil
	}

	return &game
}
