package main

import (
	"io"
	"log"
	"os"

	"github.com/godcong/excavator"
	"github.com/godcong/excavator/db"
)

func main() {
	logSet()
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//go db.PoolGetInsert(ctx)

	root := excavator.NewRoot("http://tool.httpcn.com", "/KangXi/BuShou.html")
	root.Self()
	root.IteratorSelf(func(radical *excavator.Radical) error {
		log.Println(radical)
		err := db.RD().Insert(&radical.RadicalCharacter)
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	//root.Iterator(func(radical *excavator.Radical) error {
	//
	//	for _, v := range radical.SelfCharacters() {
	//		db.RD().Insert(v)
	//	}
	//	return nil
	//})

	log.Println("wait for done")
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
