package examples

import (
	"fmt"
	"log"

	"github.com/globalsign/mgo/bson"
	"github.com/godcong/excavator"
	"github.com/godcong/excavator/db"
)

//ExampleRadical how to get the radical
func ExampleRadical() {
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
}
