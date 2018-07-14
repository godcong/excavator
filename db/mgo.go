package db

import (
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

func init() {
	session = Dial()
}

func DB() *mgo.Collection {
	return session.DB("fate").C("data")
}

func Dial() *mgo.Session {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	credential := &mgo.Credential{
		Username:    "root",
		Password:    "v2RgzSuIaBlx",

	}
	session.Login(credential)
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session
}

func Close() {
	session.Close()
}
