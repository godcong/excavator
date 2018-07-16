package main

import (
	"context"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/godcong/excavator"
	"github.com/godcong/excavator/db"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	max := 10
	logSet()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go db.PoolInsertLoop(ctx)

	var rcs []excavator.RadicalCharacter
	err := db.DB("radical").Find(bson.M{}).All(&rcs)
	if err != nil {

	}

	root := excavator.NewRoot("http://tool.httpcn.com", "/KangXi/BuShou.html")
	var radical excavator.Radical
	radical.SetRoot(root)
	ch := make(chan int, max)
	size := len(rcs)
	idx := 0
	for i := 0; i < max; i++ {
		//log.Println("thread:", idx)
		go threadLoop(idx, &radical, &rcs[idx], ch)
		idx++
	}
	for {
		if idx >= size {
			break
		}

		select {
		case <-ch:
			log.Println("thread:", idx)
			go threadLoop(idx, &radical, &rcs[idx], ch)
			idx++
		default:

		}

	}

}

func threadLoop(idx int, radical *excavator.Radical, rc *excavator.RadicalCharacter, ch chan<- int) {
	c := radical.Character(rc)
	if c.Character == "" {
		db.InsertIfNotExist("radicalwrong", rc)
		ch <- idx
		return
	}
	err := db.InsertIfNotExist("character", c)
	if err != nil {
		ch <- -1
		log.Println(err)
		return
	}
	ch <- -1
}

func logSet() {
	log.SetFlags(log.Lshortfile)
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_APPEND|os.O_SYNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	out := io.MultiWriter(f, os.Stdout)
	log.SetOutput(out)
}
