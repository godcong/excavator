package main

import (
	"log"

	"github.com/godcong/excavator"
	"github.com/godcong/excavator/db"
)

func main() {
	log.SetFlags(log.Llongfile)
	root := excavator.NewRoot("http://tool.httpcn.com","/KangXi/BuShou.html")
	root.Self()
	root.Iterator(func(radical *excavator.Radical) error {
		radical.Iterator(func(character *excavator.Character) error {
			db.DB().Insert(character)
			log.Println("insert:", *character)
			return nil
		})
		return nil
	})
	//root.GetList()
}
