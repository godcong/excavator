package examples

import (
	"github.com/godcong/fate"
)

//ExampleCharacter to add characters
func ExampleCharacter() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	//max := 10
	//
	//var rcs []excavator.RadicalCharacter
	//err := db.DB("radical").Find(bson.M{}).All(&rcs)
	//if err != nil {
	//
	//}
	//root := excavator.NewRoot("http://tool.httpcn.com", "/KangXi/BuShou.html")
	//var radical excavator.Radical
	//radical.SetRoot(root)
	//ch := make(chan bool, max)
	//size := len(rcs)
	//idx := 0
	//for i := 0; i < max; i++ {
	//	//log.Println("thread:", idx)
	//	go threadLoop(&radical, &rcs[idx], ch)
	//	idx++
	//}
	//for {
	//	if idx >= size {
	//		break
	//	}
	//
	//	select {
	//	case v := <-ch:
	//		if !v {
	//			log.Println("wrong id:", idx, rcs[idx])
	//		}
	//		//log.Println("thread:", idx)
	//		go threadLoop(&radical, &rcs[idx], ch)
	//		idx++
	//	default:
	//
	//	}
	//
	//}
}

// RadicalUnion ...
type RadicalUnion struct {
	String                *string
	RadicalCharacterArray []*fate.Character
}

func threadLoop(radical *[]RadicalUnion, rc *fate.Character, ch chan<- bool) {
	//c := radical.Character(rc)
	//if c.Character == "" {
	//	ch <- false
	//	return
	//}
	//err := db.InsertOrUpdate("character", c)
	//if err != nil {
	//	ch <- false
	//	log.Println(err)
	//	return
	//}
	//ch <- true
}
