package main

import (
	"io"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/godcong/excavator"
	"github.com/godcong/excavator/db"
)

func main() {
	max := runtime.NumCPU() * 2
	runtime.GOMAXPROCS(max)

	logSet()
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

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
		case v := <-ch:
			if v != -1 {
				log.Println("wrong id:", idx, rcs[idx])
			}
			if idx <= 20891 {
				go threadLoop(idx, &radical, &rcs[idx], ch)
				time.Sleep(3 * time.Second)
				idx++
			}
		case <-time.After(10 * time.Second):
			break
		default:

		}
	}

	//db.PoolInsertLoop(ctx)

}

func threadLoop(idx int, radical *excavator.Radical, rc *excavator.RadicalCharacter, ch chan<- int) {
	c := radical.Character(rc)
	if c.Character == "" {
		log.Println(*c)
		ch <- idx
		return
	}
	err := db.DB("character").Insert(c)
	//InsertIfNotExist("character", c)

	//db.PoolInsertAdd(c)
	if err != nil {
		ch <- idx
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
