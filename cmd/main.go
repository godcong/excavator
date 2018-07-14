package main

import (
	"flag"
	"log"

	"github.com/godcong/excavator"
	"github.com/godcong/excavator/db"
)

func main() {

	flag.Parse()
	log.SetFlags(log.Llongfile)
	root := excavator.Self()
	root.ListProcess(func(c *excavator.Character) {
		db.DB().Insert(c)
		log.Println("insert:", * c)
	})
	//root.GetList()
}
