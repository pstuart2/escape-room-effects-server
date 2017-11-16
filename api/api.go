package api

import "gopkg.in/mgo.v2"

type ErrorResponse struct {
	Message string `json:"message"`
}

type Server struct {
	Db *mgo.Session
}

func (s Server) getDb() *mgo.Session {
	return s.Db.Copy()
}

func getGameCollection(db *mgo.Session) *mgo.Collection {
	return db.DB("").C("game")
}
