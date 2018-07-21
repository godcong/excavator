package db

import "gopkg.in/mgo.v2/bson"

type IteratorFunc func(v interface{}) error

type WuXing struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	WuXing  []string      `bson:"wu_xing" json:"wu_xing"`
	Fortune string        `bson:"fortune" json:"fortune"`
	Comment string        `bson:"comment" json:"comment"`
}

type DaYan struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Index   int           `bson:"index" json:"index"`
	Sex     string        `bson:"sex"`
	Fortune string        `bson:"fortune" json:"fortune"`
	TianJiu string        `bson:"tian_jiu" json:"tian_jiu"`
	Comment string        `bson:"comment" json:"comment"`
}

type Iterable interface {
	HasNext() bool
	Next() interface{}
	Reset()
	Add(v interface{})
	Size() int
	IteratorFunc(f IteratorFunc) error
}
