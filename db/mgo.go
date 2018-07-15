package db

import (
	"context"
	"sync"
	"time"

	"github.com/godcong/excavator"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

func init() {
	session = Dial()
}

var pool sync.Pool

func DB() *mgo.Collection {
	return session.DB("fate").C("data")
}

func RD() *mgo.Collection {
	return session.DB("fate").C("radical")
}


func Dial() *mgo.Session {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	credential := &mgo.Credential{
		Username: "root",
		Password: "v2RgzSuIaBlx",
	}
	session.Login(credential)
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session
}

func Close() {
	session.Close()
}

func Insert(v interface{}) {
	pool.Put(v)
}

func PoolGetInsert(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for {
				if v := pool.Get(); v != nil {
					DB().Insert(v.(*excavator.Character))
				} else {
					return
				}
			}
		default:
			if v := pool.Get(); v != nil {
				DB().Insert(v.(*excavator.Character))
				continue
			}
			time.Sleep(10 * time.Second)
		}
	}
}
