package db

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/godcong/excavator"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

func init() {
	session = Dial()
}

var pool sync.Pool
var collections = make(map[string]*mgo.Collection)

func DB(cname string) *mgo.Collection {
	if v, b := collections[cname]; b {
		return v
	}
	collections[cname] = session.DB("fate").C(cname)
	return collections[cname]
}

func Dial() *mgo.Session {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	credential := &mgo.Credential{
		Username: "root",
		Password: "v2RgzSuIaBlx",
	}
	session.Login(credential)
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session
}

func InsertIfNotExist(name string, v interface{}) error {
	count, err := DB(name).Find(v).Count()
	if err != nil || count != 0 {
		return err
	}

	err = DB(name).Insert(v)
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	session.Close()
}

func PoolInsertAdd(v interface{}) {
	pool.Put(v)
}

func PoolInsertLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for {
				if v := pool.Get(); v != nil {
					log.Println("insert")
					err := InsertIfNotExist("character", v.(*excavator.Character))
					log.Println(err)
				} else {
					return
				}
			}
		default:
			if v := pool.Get(); v != nil {
				log.Println("insert")
				err := InsertIfNotExist("character", v.(*excavator.Character))
				log.Println(err)
				//DB("character").Insert(v.(*excavator.Character))
				continue
			}
			time.Sleep(10 * time.Second)
		}
	}
}

func InsertRootFromJson(name string, db string) {
	var rcs []*excavator.RootCharacter
	file, err := os.OpenFile(name, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	dec := json.NewDecoder(r)
	err = dec.Decode(&rcs)
	if err != nil {
		panic(err)
	}
	log.Println("size:", len(rcs))
	for idx := range rcs {
		InsertIfNotExist(db, &rcs[idx])
	}

}
func InsertRadicalFromJson(name string, db string) {
	var rcs []*excavator.RadicalCharacter
	file, err := os.OpenFile(name, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	dec := json.NewDecoder(r)
	err = dec.Decode(&rcs)
	if err != nil {
		panic(err)
	}

	log.Println("size:", len(rcs))
	for idx := range rcs {
		InsertIfNotExist(db, &rcs[idx])
	}

}

func InsertWuXingFromJson(name string, db string) {
	var wxs []*WuXing
	file, err := os.OpenFile(name, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	dec := json.NewDecoder(r)
	err = dec.Decode(&wxs)
	if err != nil {
		panic(err)
	}

	log.Println("size:", len(wxs))
	for idx := range wxs {
		InsertIfNotExist(db, &wxs[idx])
	}

}

func InsertDaYanFromJson(name string, db string) {
	var objs []*DaYan
	file, err := os.OpenFile(name, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	dec := json.NewDecoder(r)
	err = dec.Decode(&objs)
	if err != nil {
		panic(err)
	}

	log.Println("size:", len(objs))
	for idx := range objs {
		InsertIfNotExist(db, &objs[idx])
	}

}

func InsertFromJson(file string, v interface{}) {
	//TODO:
}

func getInterfaceName(v interface{}) string {
	vo := reflect.ValueOf(v)
	t := vo.Type().String()

	if idx := strings.LastIndex(t, "."); idx >= 0 && len(t) > 1 {
		return t[idx+1:]
	}
	return ""
}

func InsertCharacterFromJson(name string, db string) {
	var objs []*Character
	file, err := os.OpenFile(name, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	r := bufio.NewReader(file)
	dec := json.NewDecoder(r)
	err = dec.Decode(&objs)
	if err != nil {
		panic(err)
	}

	log.Println("size:", len(objs))
	for idx := range objs {
		DB(db).Insert(&objs[idx])
	}
}

//UpdateCommonly UpdateCommonly
func UpdateCommonly(v *excavator.CommonlyCharacter) {
	err := DB("commonly").Insert(v)
	if err != nil {
		log.Println(err)
		return
	}
}
