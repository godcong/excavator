package excavator_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/godcong/excavator"
	"github.com/godcong/excavator/db"
	"gopkg.in/mgo.v2/bson"
)

func TestRadical_Add(t *testing.T) {
	text := "汉字五行：土　是否为常用字：否"
	s := strings.SplitAfter(text, "：")
	log.Println(s, len(s))
}

func TestRoot_IteratorFunc(t *testing.T) {
	root := excavator.NewRoot("http://tool.httpcn.com", "/KangXi/BuShou.html")
	root.Self()
	root.SetBefore(func(rc *excavator.RootCharacter) error {
		count, err := db.DB("root").Find(rc).Count()
		if err != nil || count != 0 {
			return nil
		}
		db.DB("root").Insert(rc)
		return nil
	})

	radicals := root.IteratorFunc(func(rc *excavator.RootCharacter) error {
		return nil
	})
	t.Log(len(radicals))
}

func TestRadical_IteratorFunc(t *testing.T) {
	var rcs []excavator.RootCharacter
	err := db.DB("root").Find(bson.M{}).All(&rcs)
	if err != nil {
		t.Error(err)
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
