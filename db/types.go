package db

import "gopkg.in/mgo.v2/bson"

type IteratorFunc func(v interface{}) error

type WuXing struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	WuXing  []string      `bson:"wu_xing"`
	Fortune string        `bson:"fortune"`
	Comment string        `bson:"comment"`
}

type Iterable interface {
	HasNext() bool
	Next() interface{}
	Reset()
	Add(v interface{})
	Size() int
	IteratorFunc(f IteratorFunc) error
}