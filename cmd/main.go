package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/godcong/excavator"
	"github.com/godcong/excavator/db"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	logSet()
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//go db.PoolGetInsert(ctx)

	var rcs []excavator.RootCharacter
	err := db.DB("root").Find(bson.M{}).All(&rcs)
	if err != nil {

	}
	root := excavator.NewRoot("http://tool.httpcn.com", "/KangXi/BuShou.html")

	for idx := range rcs {
		radical := root.Radical(&rcs[idx])
		radical.SetRoot(root)
		radical.SetBefore(func(rc *excavator.RadicalCharacter) error {
			count, err := db.DB("radical").Find(rc).Count()
			if err != nil || count != 0 {
				log.Println(err, count)
				return fmt.Errorf("%s", "data found")
			}
			err = db.DB("radical").Insert(rc)
			if err != nil {
				log.Println(err)
			}
			return nil
		})
		radical.IteratorFunc(func(rc *excavator.RadicalCharacter) error {
			return nil
		})
	}
	//root.WaitForDone()

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
